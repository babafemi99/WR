package rest

import (
	"context"
	"errors"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func (a *API) DoPersistStaff(staff model.Staff) (*model.Staff, string, string, error) {
	// check if email already exists
	exist, err := a.Deps.Repository.EmailExist(staff.Email, staff.Role)
	if exist {
		return nil, values.Conflict, "staff with this email already exist", errors.New("duplicate resource")
	}
	if err != nil {
		return nil, values.Error, "unable to fetch user details", err
	}

	staff.Id = ulid.Make()
	staff.CreatedAt = time.Now()
	staff.UpdatedAt = time.Now()
	staff.Status = values.UserDefaultStatus

	var message, status string
	err = a.Deps.Repository.RunInTx(context.Background(), func() error {

		// generate temporary password
		pass := util.GenerateTempPassword()
		staff.HashPassword, err = util.HashPassword([]byte(pass))
		if err != nil {
			message = "unable to generate temporary password"
			status = values.Error
			return err
		}

		err = a.Deps.Repository.PersistStaff(staff)
		if err != nil {
			message = "unable add new staff"
			status = values.Error
			return err
		}

		// send email to user

		return nil
	})
	if err != nil {
		return nil, status, message, err
	}
	return &staff, values.Success, "staff added successfully", nil
}

func (a *API) DoStaffLogin(req model.LoginReq) (*model.AuthStaff, string, string, error) {
	// validate req

	// fetch user with that email
	user, err := a.Deps.Repository.FindStaffByEmail(req.Email)
	if err != nil {
		log.Println(err, "error")
		return nil, values.NotAuthorised, "Invalid credentials", err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password))
	if err != nil {
		return nil, values.NotAuthorised, "Invalid credentials", err
	}

	// createToken
	_, token, status, message, err := a.CreateAuthToken(user.Email, user.Id.String(), user.Role)
	if err != nil {
		return nil, status, message, err
	}

	// return response
	return &model.AuthStaff{
		Staff: &user,
		Auth: model.TokenInfo{
			Token:        token[0],
			RefreshToken: token[1],
		},
	}, values.Success, "log in successful", nil
}

func (a *API) DoBlockStaff(id string) (string, string, error) {

	exist, err := a.Deps.Repository.IdExist(id, "staff")
	if err != nil {
		return "", "", err
	}

	if !exist {
		return values.NotFound, "no staff with this id ", errors.New("invalid resource")
	}
	if err != nil {
		return values.Error, "unable to fetch user details", err
	}

	err = a.Deps.Repository.BLockStaff(id)
	if err != nil {
		return values.Error, "unable to block staff", err
	}

	return values.Success, "blocked staff successfully", nil
}

func (a *API) DeleteStaff(id string) (string, string, error) {
	err := a.Deps.Repository.DeleteAdmin(id)
	if err != nil {
		return values.Error, "unable to delete staff", err
	}
	return values.Success, "blocked staff successfully", nil
}

// change password -- staff

func (a *API) ChangeStaffPassword(req model.ChangePasswordReq) (string, string, error) {
	// verify ChangePasswordReq

	// load staff details

	staff, err := a.Deps.Repository.FindStaffById(req.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return values.NotFound, "no staff with this id", errors.New("invalid resource")
		}
		return values.Error, "failed to find staff", err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(staff.HashPassword), []byte(req.OldPassword))
	if err != nil {
		return values.NotAuthorised, "Invalid credentials", err
	}

	// store new password

	err = a.Deps.Repository.UpdateStaffPassword(req.UserID, req.NewPassword)
	if err != nil {
		return values.Error, "failed to update password", err
	}

	return values.Success, "updated password successfully", nil

}

func (a *API) UpdateStaffPassword(req model.UpdatePasswordReq) (string, string, error) {

	err := a.Deps.Repository.UpdateStaffPassword(req.UserID, req.Password)
	if err != nil {
		return values.Error, "failed to update password", err
	}

	return values.Success, "updated password successfully", nil
}
