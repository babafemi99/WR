package rest

import (
	"encoding/json"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *API) StaffRoutes() chi.Router {
	r := chi.NewRouter()

	r.Method(http.MethodPost, "settings/staff/change-password/{id}", Handler(a.ChangePasswordStaff))
	r.Method(http.MethodPost, "/live/{id}", Handler(a.ToggleWeddingLive))
	r.Method(http.MethodPost, "/offline/{id}", Handler(a.OffWedding))

	return r
}

func (a *API) ChangePasswordStaff(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get req from body and validate
	var newReq model.ChangePasswordReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}
	status, message, err := a.ChangeStaffPassword(newReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Err:        nil,
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (a *API) ToggleWeddingLive(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get link from url
	wID := chi.URLParam(r, "wid")
	// verify if link is for today and has not been toggled
	_, status, message, err := a.VerifyWeddingId(wID)
	if err != nil {
		return respondWithError(err, message, status)
	}

	// get req from body and validate
	var toggleReq model.ToggleWeddingReq
	err = json.NewDecoder(r.Body).Decode(&toggleReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	// set link status to toggled with user details
	status, message, err = a.ToggleWedding(toggleReq)
	if err != nil {
		return respondWithError(err, status, message)
	}

	return &ServerResponse{
		Err:        nil,
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    nil,
	}
}
func (a *API) OffWedding(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get link from url
	wID := chi.URLParam(r, "wid")

	// set link status to toggled with user details
	status, message, err := a.ToggleWeddingOff(wID)
	if err != nil {
		return respondWithError(err, status, message)
	}

	return &ServerResponse{
		Err:        nil,
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    nil,
	}
}
