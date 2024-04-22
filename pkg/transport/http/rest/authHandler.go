package rest

import (
	"encoding/json"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *API) SuperAdminRoutes() chi.Router {
	r := chi.NewRouter()
	r.Method(http.MethodPost, "/sign-in", Handler(a.SuperAdminSignIn))
	return r
}

func (a *API) AdminAuthRoutes() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "/sign-in", Handler(a.AdminSignIn))
	r.Method(http.MethodPost, "/refresh-token", Handler(a.AdminRefreshToken))
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

	admin, status, message, err := a.DoAdminLogin(r.Context(), newReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    admin,
	}
}

func (a *API) AdminRefreshToken(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	var newReq model.RefreshToken
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid refresh token", values.BadRequestBody)
	}

	res, status, message, err := a.DoAdminRefreshToken(newReq.RefreshToken)
	if err != nil {
		return respondWithError(err, message, status)
	}
	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    res,
	}
}

func (a *API) StaffSignIn(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	var newReq model.LoginReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	admin, status, message, err := a.DoStaffLogin(r.Context(), newReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Err:        nil,
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    admin,
	}
}

func (a *API) SuperAdminSignIn(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	var newReq model.LoginReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	admin, status, message, err := a.DoSuperAdminLogin(newReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    admin,
	}
}
