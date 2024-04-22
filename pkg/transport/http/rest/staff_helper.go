package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

func (a *API) DoPersistStaff(ctx context.Context, staff model.Staff) (*model.PersistRes, string, string, error) {
	// check if email already exists
	exist, err := a.Deps.Repository.EmailExist(ctx, staff.Email, "staff")
	if exist {
		return nil, values.Conflict, "staff with this email already exist", errors.New("duplicate resource")
	}
	if err != nil {
		log.Println("error", err)
		return nil, values.Error, "unable to fetch user details", err
	}

	staff.Id = uuid.New()
	staff.CreatedAt = time.Now()
	staff.UpdatedAt = time.Now()
	staff.Status = values.UserDefaultStatus
	staff.Role = "staff"

	var message, status string
	var res model.PersistRes
	err = a.Deps.Repository.RunInTx(ctx, func() error {

		// generate temporary password
		pass := util.GenerateTempPassword()
		staff.HashPassword, err = util.HashPassword([]byte(pass))
		if err != nil {
			message = "unable to generate temporary password"
			status = values.Error
			return err
		}

		err = a.Deps.Repository.PersistStaff(ctx, staff)
		if err != nil {
			message = "unable add new staff"
			status = values.Error
			return err
		}

		res.Email = staff.Email
		res.Password = pass

		// send email to user
		data := struct {
			Type     string
			Email    string
			Password string
		}{
			Type:     "User",
			Email:    res.Email,
			Password: res.Password,
		}
		patterns := []string{"welcome_user.tmpl"}
		err = a.Deps.IMailer.SendEmail("ooluwa27@gmail.com", nil, data, patterns...)
		if err != nil {
			message = "failed to send email"
			status = values.Failed
			return err
		}

		return nil
	})
	if err != nil {
		return nil, status, message, err
	}
	return &res, values.Success, "staff added successfully", nil
}

func (a *API) DoStaffLogin(ctx context.Context, req model.LoginReq) (*model.AuthStaff, string, string, error) {
	// validate req

	// fetch user with that email
	user, err := a.Deps.Repository.FindStaffByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, values.NotAuthorised, "Invalid credentials", err
		}
		return nil, values.Failed, "system error", err
	}

	if strings.ToLower(user.Status) == "blocked" {
		return nil, values.NotAuthorised, "you have been blocked", errors.New("blocked")
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

func (a *API) DoBlockStaff(ctx context.Context, id string) (string, string, error) {

	exist, err := a.Deps.Repository.IdExist(ctx, id, "staff")
	if err != nil {
		return values.Failed, "failed to find by ID", err
	}

	if !exist {
		return values.NotFound, "no staff with this id ", errors.New("invalid resource")
	}
	if err != nil {
		return values.Error, "unable to fetch user details", err
	}

	var status, message string
	err = a.Deps.Repository.RunInTx(ctx, func() error {

		err = a.Deps.Repository.BLockStaff(ctx, id)
		if err != nil {
			status = values.Failed
			message = "unable to block staff"
			return err
		}

		err = a.Deps.Redis.DeleteAuthSession(ctx, []string{fmt.Sprintf("staff-%s", id)})
		if err != nil {
			status = values.Error
			message = "failed to delete staff session"
			return err
		}
		return nil
	})
	if err != nil {
		return status, message, err
	}

	return values.Success, "blocked staff successfully", nil
}

func (a *API) DeleteStaff(ctx context.Context, id string) (string, string, error) {
	log.Println("idd", id)
	err := a.Deps.Repository.DeleteStaff(ctx, id)
	if err != nil {
		return values.Error, "unable to delete staff", err
	}
	return values.Success, "deleted staff successfully", nil
}

// change password -- staff

func (a *API) ChangeStaffPassword(ctx context.Context, req model.ChangePasswordReq) (string, string, error) {
	// verify ChangePasswordReq

	// get executor request
	executor, ok := ctx.Value(values.Executor).(model.Executor)
	if !ok {
		return values.Failed, "system error", errors.New("failed to gex executor")
	}

	// load staff details

	staff, err := a.Deps.Repository.FindStaffById(ctx, executor.Id)
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

	//encrypt new password
	hashPassword, err := util.HashPassword([]byte(req.NewPassword))
	if err != nil {
		return values.Error, "failed to hash new password", err
	}

	// store new password
	err = a.Deps.Repository.UpdateStaffPassword(ctx, executor.Id, hashPassword)
	if err != nil {
		return values.Error, "failed to update password", err
	}

	return values.Success, "updated password successfully", nil

}

func (a *API) UpdateStaffPassword(ctx context.Context, req model.UpdatePasswordReq) (string, string, error) {

	//encrypt new password
	hashPassword, err := util.HashPassword([]byte(req.Password))
	if err != nil {
		return values.Error, "failed to hash new password", err
	}

	err = a.Deps.Repository.UpdateStaffPassword(ctx, req.UserID, hashPassword)
	if err != nil {
		return values.Error, "failed to update password", err
	}

	return values.Success, "updated password successfully", nil
}
