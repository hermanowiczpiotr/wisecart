package server

import (
	"github.com/go-chi/render"
	"net/http"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	httpStatus int
}

func BadRequestError(err error, w http.ResponseWriter, r *http.Request) {
	resp := ErrorResponse{err.Error(), http.StatusBadRequest}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
