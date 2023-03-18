package shopify

import (
	"encoding/json"
	goshopify "github.com/bold-commerce/go-shopify"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/application/dto"
	"github.com/hermanowiczpiotr/wisecart/internal/cart/domain/entity"
)

type authData struct {
	ApiKey    string `json:"api_key"`
	SecretKey string `json:"secret_key"`
}

type Client struct {
}

func (c *Client) FetchProducts(storeProfile *entity.StoreProfile) (dto.ProductDtoList, error) {

	authData, err := parseAuthData(storeProfile)

	if err != nil {
		return dto.ProductDtoList{}, err
	}

	app := goshopify.App{
		ApiKey:    authData.ApiKey,
		ApiSecret: authData.SecretKey,
	}

	client := goshopify.NewClient(app, storeProfile.Name, "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")
	shopifyProducts, err := client.Product.List(nil)

	if err != nil {
		return dto.ProductDtoList{}, err
	}

	productDtoList := make([]dto.ProductDto, len(shopifyProducts))

	for i, shopifyProduct := range shopifyProducts {
		productDtoList[i] = dto.ProductDto{
			ShopifyId:   shopifyProduct.ID,
			Title:       shopifyProduct.Title,
			Description: shopifyProduct.BodyHTML,
		}
	}

	return dto.ProductDtoList{
		Products: productDtoList,
	}, nil
}

func (c *Client) Support(storeProfile *entity.StoreProfile) bool {
	return storeProfile.Type == "shopify"
}

func parseAuthData(storeProfile *entity.StoreProfile) (authData, error) {
	var data authData
	err := json.Unmarshal(storeProfile.AuthorizationData, &data)

	if err != nil {
		return authData{}, err
	}

	return data, nil
}
