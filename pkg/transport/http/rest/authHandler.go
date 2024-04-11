package rest

import (
	"encoding/json"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *API) AdminAuthRoutes() chi.Router {
	r := chi.NewRouter()
	//r.Use(RequestTracing)

	r.Method(http.MethodPost, "/sign-in", Handler(a.AdminSignIn))
	return r
}

func (a *API) StaffAuthRoutes() chi.Router {
	r := chi.NewRouter()
	r.Method(http.MethodPost, "/sign-in", Handler(a.StaffSignIn))
	return r
}

func (a *API) ExternalAuthRoutes() chi.Router {
	r := chi.NewRouter()
	r.Method(http.MethodPost, "/sign-in", Handler(a.AdminSignIn))
	return r
}

func (a *API) AdminSignIn(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	var newReq model.LoginReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	admin, status, message, err := a.DoAdminLogin(newReq)
	if err != nil {
		respondWithError(err, message, status)
	}

	return &ServerResponse{
		Err:        nil,
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    admin,
	}
}

func (a *API) StaffSignIn(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	var newReq model.LoginReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	admin, status, message, err := a.DoStaffLogin(newReq)
	if err != nil {
		respondWithError(err, message, status)
	}

	return &ServerResponse{
		Err:        nil,
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    admin,
	}
}
