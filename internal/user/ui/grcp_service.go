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
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

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
		return nil, err
	}

	qry := query.GetUserByEmailQuery{
		Email: registerPayload.Email,
	}

	user, _ := gs.application.GetUserByEmailQueryHandler.Handle(qry)

	if user != nil {
		return &genproto.RegisterResponse{
			Status: http.StatusConflict,
			Error:  "Email is already in use",
		}, nil
	}

	cmd := command.AddUserCommand{
		Email:    registerPayload.Email,
		Password: string(hashedPassword),
	}

	err = gs.application.AddUserCommandHandler.Handle(cmd)

	return &genproto.RegisterResponse{
		Status: http.StatusAccepted,
	}, err
}

func (gs GRPCService) Login(ctx context.Context, loginPayload *genproto.LoginRequest) (*genproto.LoginResponse, error) {
	qry := query.GetUserByEmailQuery{loginPayload.Email}

	user, _ := gs.application.GetUserByEmailQueryHandler.Handle(qry)
	if user == nil {
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

	jwtauth.SetExpiry(claims, time.Now().Add(TOKEN_EXPIRATION_TIME))
	_, tokenString, err := gs.tokenAuth.Encode(claims)
	if err != nil {
		return nil, err
	}

	return &genproto.LoginResponse{
		Status: http.StatusOK,
		Token:  tokenString,
	}, nil
}

func (gs GRPCService) Validate(ctx context.Context, validateRequest *genproto.ValidateRequest) (*genproto.ValidateResponse, error) {
	jwtToken, err := jwtauth.VerifyToken(gs.tokenAuth, validateRequest.Token)
	if err != nil {
		return nil, err
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
		UserId: userId.(int64),
	}, nil
}
