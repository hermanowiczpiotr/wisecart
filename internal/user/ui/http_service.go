package ui

import (
	"errors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt"
	"github.com/hermanowiczpiotr/ola/internal/user/application"
	"github.com/hermanowiczpiotr/ola/internal/user/application/command"
	"github.com/hermanowiczpiotr/ola/internal/user/application/query"
	"github.com/hermanowiczpiotr/ola/internal/user/infrastructure/server"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var TOKEN_EXPIRATION_TIME time.Duration = time.Minute * 15

type UsersHandler struct {
	usersApp  application.UserApp
	tokenAuth *jwtauth.JWTAuth
}

func NewUsersHandler(ua application.UserApp, ta *jwtauth.JWTAuth) UsersHandler {
	return UsersHandler{
		usersApp:  ua,
		tokenAuth: ta,
	}
}

func (uh UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginPayload := Login{}
	err := render.Decode(r, &loginPayload)
	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	qry := query.GetUserByEmailQuery{
		Email: loginPayload.Email,
	}

	user, err := uh.usersApp.GetUserByEmailQueryHandler.Handle(qry)
	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password))
	if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
		server.BadRequestError(errors.New("unauthorized"), w, r)
		return
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	}

	jwtauth.SetExpiry(claims, time.Now().Add(TOKEN_EXPIRATION_TIME))
	_, tokenString, err := uh.tokenAuth.Encode(claims)

	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	render.Respond(w, r, tokenString)
}

func (uh UsersHandler) GetUserById(w http.ResponseWriter, r *http.Request, userid string) {
	qry := query.GetUserByIdQuery{
		userid,
	}

	user, err := uh.usersApp.GetUserByIdQueryHandler.Handle(qry)

	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	render.Respond(w, r, user)
}

func (uh UsersHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	newUser := NewUser{}

	err := render.Decode(r, &newUser)

	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)

	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	qry := query.GetUserByEmailQuery{newUser.Email}
	user, _ := uh.usersApp.GetUserByEmailQueryHandler.Handle(qry)

	if user != nil {
		server.BadRequestError(errors.New("email is already taken"), w, r)
		return
	}

	cmd := command.AddUserCommand{
		Email:    newUser.Email,
		Password: string(hashedPassword),
	}

	err = uh.usersApp.AddUserCommandHandler.Handle(cmd)

	if err != nil {
		server.BadRequestError(err, w, r)
		return
	}

	render.Respond(w, r, User{
		Email: newUser.Email,
	})
}

func (uh UsersHandler) Verifier() func(http.Handler) http.Handler {
	return jwtauth.Verifier(uh.tokenAuth)
}
