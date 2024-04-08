package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/babafemi99/WR/internal/config"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/deps"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"time"
)

type Handler func(w http.ResponseWriter, r *http.Request) *ServerResponse

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response := h(w, r)
	responseByte, err := json.Marshal(response)
	if err != nil {
		writeErrorResponse(w, err, values.Error, "unable to marshal server response")
		return
	}
	writeJSONResponse(w, responseByte, response.StatusCode)
}

type API struct {
	Server *http.Server
	Config *config.Config
	Deps   *deps.Dependencies
}

// Serve starts the core service
func (a *API) Serve() error {
	a.Server = &http.Server{
		Addr:           fmt.Sprintf(":%d", a.Config.Port),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        a.setUpServerHandler(),
		MaxHeaderBytes: 1024 * 1024,
	}

	return a.Server.ListenAndServe()
}

func (a *API) setUpServerHandler() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodPost, http.MethodGet, http.MethodPatch, http.MethodPut, http.MethodDelete},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-For", "Content-Type", "X-CSRF-Token", values.HeaderRequestID,
			values.HeaderRequestSource},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	return mux
}

// Shutdown shuts the core service down
func (a *API) Shutdown() error {
	// todo (ore): shut down database

	err := a.Server.Shutdown(context.Background())
	if err != nil {
		return err
	}
	return nil
}
