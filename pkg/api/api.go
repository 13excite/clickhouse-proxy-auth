// Package api implements the main clickhouse auth api
package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

// Service are the methods used to retrieve the API data
type Service interface {
	AuthClickhouseReq(context.Context) (AuthResponse, error)
}

// handler implements the main api logic
type handler struct {
	service Service
}

// AuthResponse is the response for a authClickhouse request
type AuthResponse struct {
	Status string `json:"status"`
}

// HandlerOpt is an option func for the handler
type HandlerOpt func(*handler)

// NewHandler creates a new Clickhouse Auth API
func NewHandler(service Service, opts ...HandlerOpt) http.Handler {
	h := &handler{
		service: service,
	}
	for _, o := range opts {
		o(h)
	}

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/auth", h.authClickhouse)
	})
	return router
}

// TODO: add logic
func (h *handler) authClickhouse(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "OK"})

}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
