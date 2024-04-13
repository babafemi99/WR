package rest

import (
	"encoding/json"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (a *API) AdminRoutes() chi.Router {
	r := chi.NewMux()
	adminRoutes := r.With(AuthenticateAdmin)
	superAdminRoutes := r.With(AuthenticateSuperAdmin)

	adminRoutes.Method(http.MethodPatch, "settings/admin/change-password/{id}", Handler(a.ChangePasswordAdmin))
	adminRoutes.Method(http.MethodPatch, "settings/staff/block/{id}", Handler(a.BlockStaff))
	adminRoutes.Method(http.MethodPost, "/create", Handler(a.LoadWeddingDetails))
	//change staff details
	// change admin details
	// load wedding details

	superAdminRoutes.Method(http.MethodPost, "/settings/admin", Handler(a.CreateAdmin))
	superAdminRoutes.Method(http.MethodPost, "/settings/staff", Handler(a.CreateStaff))
	superAdminRoutes.Method(http.MethodPatch, "settings/admin/block/{id}", Handler(a.BlockAdmin))
	superAdminRoutes.Method(http.MethodDelete, "settings/admin/delete/{id}", Handler(a.RemoveAdmin))
	superAdminRoutes.Method(http.MethodDelete, "settings/staff/delete/{id}", Handler(a.RemoveStaff))

	// this enables a super admin to change their other admin and staff password
	superAdminRoutes.Method(http.MethodPatch, "settings/staff/update-password/{id}", Handler(a.SuperModifyStaffPassword))
	superAdminRoutes.Method(http.MethodPatch, "settings/admin/update-password/{id}", Handler(a.SuperModifyAdminPassword))

	//create entities -- todo
	// block entities --todo
	//delete entities --todo

	return r
}

func (a *API) CreateAdmin(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get req from body and validate
	var newAdmin model.Admin
	err := json.NewDecoder(r.Body).Decode(&newAdmin)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	admin, status, message, err := a.DoPersistAdmin(newAdmin)
	if err != nil {
		respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    admin,
	}

}

func (a *API) CreateStaff(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get req from body and validate
	var newStaff model.Staff
	err := json.NewDecoder(r.Body).Decode(&newStaff)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}

	admin, status, message, err := a.DoPersistStaff(newStaff)
	if err != nil {
		respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
		Payload:    admin,
	}
}

func (a *API) BlockAdmin(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	Id := chi.URLParam(r, "id")
	if Id == "" {
		return respondWithError(nil, "invalid admin id", values.BadRequestBody)
	}

	status, message, err := a.DoBlockAdmin(Id)
	if err != nil {
		respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (a *API) BlockStaff(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	Id := chi.URLParam(r, "id")
	if Id == "" {
		return respondWithError(nil, "invalid staff id", values.BadRequestBody)
	}

	status, message, err := a.DoBlockStaff(Id)
	if err != nil {
		respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (a *API) RemoveAdmin(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	Id := chi.URLParam(r, "id")
	if Id == "" {
		return respondWithError(nil, "invalid admin id", values.BadRequestBody)
	}

	status, message, err := a.DeleteAdmin(Id)
	if err != nil {
		respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (a *API) RemoveStaff(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	Id := chi.URLParam(r, "id")
	if Id == "" {
		return respondWithError(nil, "invalid staff id", values.BadRequestBody)
	}

	status, message, err := a.DeleteStaff(Id)
	if err != nil {
		respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (a *API) SuperModifyStaffPassword(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get req from body and validate
	var newReq model.UpdatePasswordReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}
	status, message, err := a.UpdateStaffPassword(newReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (a *API) SuperModifyAdminPassword(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get req from body and validate
	var newReq model.UpdatePasswordReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}
	status, message, err := a.UpdateAdminPassword(newReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func (a *API) ChangePasswordAdmin(_ http.ResponseWriter, r *http.Request) *ServerResponse {
	// get req from body and validate
	var newReq model.ChangePasswordReq
	err := json.NewDecoder(r.Body).Decode(&newReq)
	if err != nil {
		return respondWithError(err, "invalid request body provided", values.BadRequestBody)
	}
	status, message, err := a.ChangeAdminPassword(newReq)
	if err != nil {
		return respondWithError(err, message, status)
	}

	return &ServerResponse{
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}
