package ui

import (
	"github.com/go-chi/render"
	"github.com/hermanowiczpiotr/ola/internal/gateway/infrastructure/server/grcp"
	"net/http"
)

type Handler struct {
	UserClient grcp.UserClient
}

func NewHandler(client grcp.UserClient) Handler {
	return Handler{
		UserClient: client,
	}
}

func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	payload := grcp.LoginRequest{}

	err := render.Decode(r, &payload)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	res, err := h.UserClient.Login(r.Context(), &payload)

	if err != nil {
		render.Respond(w, r, err)
		return
	}

	render.Respond(w, r, res)
}

func (h Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	payload := grcp.RegisterRequest{}

	err := render.Decode(r, &payload)
	if err != nil {
		render.Respond(w, r, err)
		return
	}

	res, err := h.UserClient.Register(r.Context(), &payload)

	if err != nil {
		render.Respond(w, r, res.Error)
		return
	}

	render.Respond(w, r, res)
}

func (h Handler) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			render.Respond(w, r, "Unauthorized")
			return
		}

		payload := grcp.ValidateRequest{
			Token: token,
		}

		res, err := h.UserClient.Validate(r.Context(), &payload)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			render.Respond(w, r, "Unauthorized")
			return
		}

		if res.Status != http.StatusOK {
			w.WriteHeader(http.StatusUnauthorized)
			render.Respond(w, r, "Unauthorized")
		}

		next.ServeHTTP(w, r)
	})
}
