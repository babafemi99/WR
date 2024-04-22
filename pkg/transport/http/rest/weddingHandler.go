package rest

import (
	"encoding/json"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (a *API) WeddingRoutes() chi.Router {
	weddingRouter := chi.NewRouter()

	weddingRouter.Method(http.MethodGet, "/{wid}/{code}", Handler(a.JoinWeddingMeeting))
	weddingRouter.Method(http.MethodPost, "/add-member", Handler(a.AddMember))
	weddingRouter.Method(http.MethodGet, "/{wid}/members", Handler(a.GetAllMembers))
	weddingRouter.Method(http.MethodDelete, "/{wid}/{email}", Handler(a.RemoveMember))

	return weddingRouter
}

func (a *API) JoinWeddingMeeting(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get id from url
	wID := chi.URLParam(r, "wid")
	if wID == "" {
		return respondWithError(nil, "invalid wedding id", values.BadRequestBody)
	}

	wCode := chi.URLParam(r, "code")
	if wCode == "" {
		return respondWithError(nil, "invalid wedding id", values.BadRequestBody)
	}
	// verify id

	weddingDetails, status, message, err := a.JoinWedding(r.Context(), wID, wCode)
	if err != nil {
		return respondWithError(err, message, status)
	}

	// send response
	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    weddingDetails,
	}
}

func (a *API) LoadWeddingDetails(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get req from body and validate

	var newWeddingReq model.NewWeddingReq
	err := json.NewDecoder(r.Body).Decode(&newWeddingReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	// persist in db
	wedding, status, message, err := a.DoPersistWedding(r.Context(), newWeddingReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	//return response
	return &ServerResponse{
		Status:     status,
		Message:    message,
		StatusCode: util.StatusCode(status),
		Payload:    wedding,
	}

}

func (a *API) AddMember(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	// get req from body and validate
	var newMember model.Member
	err := json.NewDecoder(r.Body).Decode(&newMember)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	status, message, err := a.DoAddMember(r.Context(), newMember)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    nil,
	}
}

func (a *API) GetAllMembers(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	wID := chi.URLParam(r, "wid")
	if wID == "" {
		return respondWithError(nil, "invalid wedding id", values.BadRequestBody)
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 20
	}

	// Calculate offset and limit based on page and pageSize
	offset := (page - 1) * pageSize
	limit := pageSize

	members, err := a.Deps.Repository.GetMembers(r.Context(), wID, offset, limit)
	if err != nil {
		return nil
	}
	//todo handle error no rows
	return &ServerResponse{
		Message:    "fetched members successfully",
		Status:     values.Success,
		StatusCode: util.StatusCode(values.Success),
		Payload:    members,
	}
}

func (a *API) RemoveMember(_ http.ResponseWriter, r *http.Request) *ServerResponse {

	wID := chi.URLParam(r, "wid")
	if wID == "" {
		return respondWithError(nil, "invalid wedding id", values.BadRequestBody)
	}

	email := chi.URLParam(r, "email")
	if wID == "" {
		return respondWithError(nil, "invalid email in url", values.BadRequestBody)
	}

	// delete the email

	err := a.Deps.Repository.RemoveMember(r.Context(), email, wID)
	if err != nil {
		return respondWithError(err, "failed to remove member", values.Error)
	}

	return &ServerResponse{
		Message:    "removed member successfully",
		Status:     values.Success,
		StatusCode: util.StatusCode(values.Success),
	}
}
