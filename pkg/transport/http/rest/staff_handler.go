package rest

import (
	"encoding/json"
	"errors"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *API) StaffRoutes() chi.Router {
	r := chi.NewRouter()
	r.Use(a.AuthenticateStaff)

	r.Method(http.MethodPatch, "/settings/staff/change-password", Handler(a.ChangePasswordStaff))
	r.Method(http.MethodPost, "/wedding/live", Handler(a.ToggleWeddingLive))
	r.Method(http.MethodPatch, "/wedding/offline/{id}", Handler(a.OffWedding))

	return r
}

func (a *API) ChangePasswordStaff(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get req from body and validate
	var newReq model.ChangePasswordReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}
	status, message, err := a.ChangeStaffPassword(r.Context(), newReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (a *API) ToggleWeddingLive(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	// get req from body and validate
	var toggleReq model.ToggleWeddingReq
	err := json.NewDecoder(r.Body).Decode(&toggleReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	// set link status to toggled with user details
	status, message, err := a.ToggleWedding(r.Context(), toggleReq)
	if err != nil {
		return respondWithError(err, status, message)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    nil,
	}
}
func (a *API) OffWedding(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get link from url
	wID := chi.URLParam(r, "id")
	if wID == "nil" {
		respondWithError(errors.New("invalid wedding id"), "bad request", values.BadRequestBody)
	}

	// set link status to toggled with user details
	status, message, err := a.ToggleWeddingOff(r.Context(), wID)
	if err != nil {
		return respondWithError(err, status, message)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    nil,
	}
}
