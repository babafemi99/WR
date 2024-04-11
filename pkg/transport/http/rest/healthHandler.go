package rest

import (
	"github.com/babafemi99/WR/internal/values"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func HealthRoutes() chi.Router {
	mux := chi.NewRouter()
	mux.Method(http.MethodGet, "/", Handler(func(w http.ResponseWriter, r *http.Request) *ServerResponse {
		return &ServerResponse{
			Message:    values.Success,
			Status:     values.Success,
			StatusCode: http.StatusOK,
		}
	}))
	return mux
}
