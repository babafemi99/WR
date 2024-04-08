package rest

import (
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *API) WeddingRoutes() chi.Router {
	weddingRouter := chi.NewRouter()

	weddingRouter.Method(http.MethodPost, "/{link}", Handler(a.JoinWeddingMeeting))

	return weddingRouter
}

func (a *API) JoinWeddingMeeting(_ http.ResponseWriter, _ *http.Request) *ServerResponse {
	return &ServerResponse{
		Err:        nil,
		Message:    "your wedding isn't live yet",
		Status:     values.NotAuthorised,
		StatusCode: util.StatusCode(values.NotAuthorised),
	}
}
