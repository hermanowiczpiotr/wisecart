package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hermanowiczpiotr/wisecart/internal/gateway/infrastructure/server/genproto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"net/http"
	"os"
)

var (
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("SHOPIFY_API_KEY"),
		ClientSecret: os.Getenv("SHOPIFY_API_SECRET"),
		RedirectURL:  fmt.Sprintf("%s", os.Getenv("APPLICATION_URL")),
		Scopes:       []string{"read_products"},
		Endpoint:     oauth2.Endpoint{},
	}
)

type Router struct {
	engine         *gin.Engine
	UserGrcpClient genproto.UserClient
	CartClient     genproto.CartClient
}

func NewRouter(userClient genproto.UserClient, cartClient genproto.CartClient) *Router {
	router := gin.Default()
	r := &Router{
		engine:         router,
		UserGrcpClient: userClient,
		CartClient:     cartClient,
	}

	router.Use(corsMiddleware())
	//router.Use(middleware.RequestID)
	//router.Use(middleware.RealIP)
	//router.Use(middleware.Logger)
	//router.Use(middleware.Recoverer)
	//router.Use(middleware.Timeout(60 * time.Second))

	router.POST("login", r.Login)
	router.POST("register", r.RegisterUser)

	router.GET("/shopify/login", r.LoginShopify)
	router.GET("/shopify/callback", r.Callback)

	apiGroup := router.Group("/api")
	apiGroup.Use(r.Validate())

	apiGroup.POST("/cart/sync_profile", r.syncProducts)

	return r
}

func (router *Router) Run(addr string) error {
	return router.engine.Run(addr)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

func (router *Router) Login(c *gin.Context) {
	payload := genproto.LoginRequest{}

	err := c.Bind(&payload)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	res, err := router.UserGrcpClient.Login(c, &payload)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	c.JSON(200, res)
}

func (router *Router) RegisterUser(c *gin.Context) {
	payload := genproto.RegisterRequest{}

	err := c.Bind(&payload)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, "Unable to register user")
		return
	}

	res, err := router.UserGrcpClient.Register(c, &payload)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, "Unable to register user")
		return
	}

	c.JSON(200, res)
}

func (router *Router) Validate() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			return
		}

		payload := genproto.ValidateRequest{
			Token: token,
		}

		res, err := router.UserGrcpClient.Validate(c, &payload)

		if err != nil {
			log.Error(err)
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			return
		}

		if res.Status != http.StatusOK {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			return
		}

		c.Next()
	}
}

func (router *Router) LoginShopify(c *gin.Context) {
	shop := c.Request.URL.Query().Get("shop")
	if shop == "" {
		c.JSON(http.StatusBadRequest, "Missing 'shop' parameter")
		return
	}

	setShopifyOAuth2Endpoints(shop)

	userId, err := router.getUserId(c)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		return
	}

	url := oauthConf.AuthCodeURL(userId)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (router *Router) Callback(c *gin.Context) {
	r := c.Request
	state := r.FormValue("state")
	code := r.FormValue("code")
	shop := r.FormValue("shop")

	setShopifyOAuth2Endpoints(shop)

	token, err := oauthConf.Exchange(c, code)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	jsonToken, _ := json.MarshalIndent(token, "", "    ")

	payload := genproto.AddStoreProfileRequest{
		UserId:            state,
		Name:              shop,
		Type:              "shopify",
		AuthorizationData: jsonToken,
	}

	_, err = router.CartClient.AddProfile(r.Context(), &payload)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusUnauthorized, "Unauthorized")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func (router *Router) syncProducts(c *gin.Context) {

	payload := genproto.SynchronizeProductsRequest{}

	err := c.Bind(&payload)

	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	response, err := router.CartClient.SynchronizeProducts(c, &payload)
	if err != nil || response.Error != "" {
		log.Error(err)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	c.Status(http.StatusAccepted)
}

func (router *Router) getUserId(c *gin.Context) (string, error) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return "", errors.New("Token not found")
	}

	payload := genproto.ValidateRequest{
		Token: token,
	}

	res, _ := router.UserGrcpClient.Validate(c, &payload)

	if len(res.Error) != 0 {
		return "", errors.New(res.Error)
	}

	return res.UserId, nil
}

func setShopifyOAuth2Endpoints(shop string) {
	oauthConf.Endpoint.AuthURL = "https://" + shop + "/admin/oauth/authorize"
	oauthConf.Endpoint.TokenURL = "https://" + shop + "/admin/oauth/access_token"
}
