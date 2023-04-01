package ui

import (
	context "context"
	"errors"
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt"
	"github.com/hermanowiczpiotr/wisecart/internal/user/application"
	"github.com/hermanowiczpiotr/wisecart/internal/user/application/command"
	"github.com/hermanowiczpiotr/wisecart/internal/user/application/query"
	"github.com/hermanowiczpiotr/wisecart/internal/user/infrastructure/genproto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

const tokenExpirationTime = 15

type GRPCService struct {
	application application.UserApp
	tokenAuth   *jwtauth.JWTAuth
}

func NewGRPCService(app application.UserApp, ta *jwtauth.JWTAuth) GRPCService {
	return GRPCService{
		application: app,
		tokenAuth:   ta,
	}
}

func (gs GRPCService) Register(ctx context.Context, registerPayload *genproto.RegisterRequest) (*genproto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerPayload.Password), 8)
	if err != nil {
		log.Error(err)
		return &genproto.RegisterResponse{
			Status: http.StatusConflict,
			Error:  err.Error(),
		}, err
	}

	qry := query.GetUserByEmailQuery{
		Email: registerPayload.Email,
	}

	user, _ := gs.application.GetUserByEmailQueryHandler.Handle(qry)

	if user != nil {
		log.Infof("email is already in use. email: %s")
		return &genproto.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "email is already in use",
		}, errors.New("email is already in use")
	}

	cmd := command.AddUserCommand{
		Email:    registerPayload.Email,
		Password: string(hashedPassword),
	}

	err = gs.application.AddUserCommandHandler.Handle(cmd)

	return &genproto.RegisterResponse{
		Status: http.StatusAccepted,
	}, nil
}

func (gs GRPCService) Login(ctx context.Context, loginPayload *genproto.LoginRequest) (*genproto.LoginResponse, error) {
	qry := query.GetUserByEmailQuery{loginPayload.Email}

	user, _ := gs.application.GetUserByEmailQueryHandler.Handle(qry)
	if user == nil {
		log.Infof("user witth email: %s not found")
		return &genproto.LoginResponse{
			Status: http.StatusNotFound,
			Error:  fmt.Sprintf("User with email: %s not found", loginPayload.Email),
		}, nil
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password))
	if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
		return &genproto.LoginResponse{
			Status: http.StatusUnauthorized,
			Error:  fmt.Sprintf("Password not match for user: %s", loginPayload.Email),
		}, nil
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
	}

	jwtauth.SetExpiry(claims, time.Now().Add(tokenExpirationTime))
	_, tokenString, err := gs.tokenAuth.Encode(claims)
	if err != nil {
		log.Info(err)
		return nil, err
	}

	return &genproto.LoginResponse{
		Status: http.StatusOK,
		Token:  tokenString,
	}, nil
}

func (gs GRPCService) Validate(ctx context.Context, validateRequest *genproto.ValidateRequest) (*genproto.ValidateResponse, error) {
	parts := strings.Split(validateRequest.Token, "Bearer ")

	if len(parts) == 2 {
		fmt.Println(parts[1])
	} else {
		return &genproto.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  "Token not found in string",
		}, nil
	}

	jwtToken, err := jwtauth.VerifyToken(gs.tokenAuth, parts[1])

	if err != nil {
		return &genproto.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  err.Error(),
		}, nil
	}

	userId, exists := jwtToken.Get("id")

	if !exists {
		return &genproto.ValidateResponse{
			Status: http.StatusUnauthorized,
			Error:  "Cannot find user",
		}, nil
	}

	return &genproto.ValidateResponse{
		Status: http.StatusOK,
		UserId: userId.(string),
	}, nil
}
