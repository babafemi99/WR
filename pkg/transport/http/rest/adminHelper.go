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

func (a *API) DoPersistAdmin(admin model.Admin) (*model.Admin, string, string, error) {
	// verify admin body

	// check if email already exists
	exist, err := a.Deps.EmailExist(admin.Email, admin.Role)
	if exist {
		return nil, values.Conflict, "admin with this email already exist", errors.New("duplicate resource")
	}
	if err != nil {
		return nil, values.Error, "unable to fetch user details", err
	}

	admin.Id = ulid.Make()
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()
	admin.Status = values.UserDefaultStatus

	var message, status string
	err = a.Deps.Repository.RunInTx(context.Background(), func() error {

		// generate temporary password
		pass := util.GenerateTempPassword()
		admin.HashPassword, err = util.HashPassword([]byte(pass))
		if err != nil {
			message = "unable to generate temporary password"
			status = values.Error
			return err
		}

		err = a.Deps.PersistAdmin(admin)
		if err != nil {
			message = "unable add new admin"
			status = values.Error
			return err
		}

		// send email to user

		return nil
	})
	if err != nil {
		return nil, status, message, err
	}

	return &admin, values.Success, "admin added successfully", nil
}

func (a *API) DoAdminLogin(req model.LoginReq) (*model.AdminAuthRes, string, string, error) {

	// validate req

	// fetch user with that email
	user, err := a.Deps.Repository.FindAdminByEmail(req.Email)
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

	// return response
	return &model.AdminAuthRes{
		Admin: &user,
		Auth: model.TokenInfo{
			Token:                  "",
			TokenExpiryTime:        time.Time{},
			RefreshToken:           "",
			RefreshTokenExpiryTime: time.Time{},
		},
	}, values.Success, "log in successful", nil
}

func (a *API) DoBlockAdmin(id string) (string, string, error) {

	exist, err := a.Deps.Repository.IdExist(id, "admin")
	if !exist {
		return values.NotFound, "no admin with this id", errors.New("invalid resource")
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

func (a *API) DeleteAdmin(id string) (string, string, error) {
	err := a.Deps.Repository.DeleteAdmin(id)
	if err != nil {
		return values.Error, "unable to delete staff", err
	}
	return values.Success, "blocked staff successfully", nil
}

// script
func (a *API) ImportWeddingDetails() (string, string, error) {
	//write script that will consumeAPI generate uniqueId for wedding and then store data in the database
	return "", "", nil
}

func (a *API) ChangeAdminPassword(req model.ChangePasswordReq) (string, string, error) {
	// verify ChangePasswordReq

	// load admin details

	staff, err := a.Deps.Repository.FindAdminById(req.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return values.NotFound, "no admin with this id", errors.New("invalid resource")
		}
		return values.Error, "failed to find admin", err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(staff.HashPassword), []byte(req.OldPassword))
	if err != nil {
		return values.NotAuthorised, "Invalid credentials", err
	}

	// store new password

	err = a.Deps.Repository.UpdateAdminPassword(req.UserID, req.NewPassword)
	if err != nil {
		return values.Error, "failed to update password", err
	}

	return values.Success, "updated password successfully", nil
}

func (a *API) UpdateAdminPassword(req model.UpdatePasswordReq) (string, string, error) {
	err := a.Deps.Repository.UpdateStaffPassword(req.UserID, req.Password)
	if err != nil {
		return values.Error, "failed to update password", err
	}

	return values.Success, "updated password successfully", nil
}
