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

func (a *API) DoPersistAdmin(ctx context.Context, admin model.Admin) (*model.PersistRes, string, string, error) {
	// verify admin body

	// check if email already exists
	exist, err := a.Deps.Repository.EmailExist(ctx, admin.Email, "admin")
	if exist {
		return nil, values.Conflict, "admin with this email already exist", errors.New("duplicate resource")
	}
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, values.Error, "system error", err
	}

	admin.Id = uuid.New()
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()
	admin.Status = values.UserDefaultStatus
	admin.Role = "admin"

	var message, status string
	var res model.PersistRes
	err = a.Deps.Repository.RunInTx(context.Background(), func() error {

		// generate temporary password
		pass := util.GenerateTempPassword()
		admin.HashPassword, err = util.HashPassword([]byte(pass))
		if err != nil {
			message = "unable to generate temporary password"
			status = values.Error
			return err
		}

		err = a.Deps.Repository.PersistAdmin(ctx, admin)
		if err != nil {
			message = "unable add new admin"
			status = values.Error
			return err
		}

		res.Email = admin.Email
		res.Password = pass

		// send email to user
		data := struct {
			Type     string
			Email    string
			Password string
		}{
			Type:     "Admin",
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

	return &res, values.Success, "admin added successfully", nil
}

func (a *API) DoAdminLogin(ctx context.Context, req model.LoginReq) (*model.AdminAuthRes, string, string, error) {
	// validate req

	// fetch user with that email
	user, err := a.Deps.Repository.FindAdminByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, values.NotAuthorised, "Invalid credentials", err
		}
		return nil, values.Failed, "system error", err
	}

	log.Println(user.Status)
	if strings.ToLower(user.Status) == "blocked" {
		return nil, values.NotAuthorised, "you have been blocked", errors.New("blocked")
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(req.Password))
	if err != nil {
		return nil, values.NotAuthorised, "Invalid credentials", err
	}

	// createToken
	_, tokenArr, status, message, err := a.CreateAuthToken(user.Email, user.Id.String(), user.Role)
	if err != nil {
		return nil, status, message, err
	}

	log.Println(tokenArr)
	// return response
	return &model.AdminAuthRes{
		Admin: &user,
		Auth: model.TokenInfo{
			Token:        tokenArr[0],
			RefreshToken: tokenArr[1],
		},
	}, values.Success, "log in successful", nil
}

func (a *API) DoBlockAdmin(ctx context.Context, id string) (string, string, error) {

	exist, err := a.Deps.Repository.IdExist(ctx, id, "admin")
	if !exist {
		return values.NotFound, "no admin with this id", errors.New("invalid resource")
	}
	if err != nil {
		return values.Error, "unable to fetch user details", err
	}

	var status, message string
	err = a.Deps.Repository.RunInTx(ctx, func() error {

		err = a.Deps.Repository.BLockAdmin(ctx, id)
		if err != nil {
			status, message = values.Error, "unable to block staff"
			return err
		}

		err = a.Deps.Redis.DeleteAuthSession(ctx, []string{fmt.Sprintf("admin-%s", id)})
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

func (a *API) DeleteAdmin(ctx context.Context, id string) (string, string, error) {

	err := a.Deps.Repository.DeleteAdmin(ctx, id)
	if err != nil {
		return values.Error, "unable to delete staff", err
	}
	return values.Success, "deleted staff successfully", nil
}

func (a *API) ImportWeddingDetails() (string, string, error) {
	//write script that will consumeAPI generate uniqueId for wedding and then store data in the database
	return "", "", nil
}

func (a *API) ChangeAdminPassword(ctx context.Context, req model.ChangePasswordReq) (string, string, error) {
	// verify ChangePasswordReq

	// getId from context
	executor, ok := ctx.Value(values.Executor).(model.Executor)
	if !ok {
		return values.Failed, "system error", errors.New("failed to get executor")
	}

	// load admin details
	staff, err := a.Deps.Repository.FindAdminById(ctx, executor.Id)
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

	//encrypt new password
	hashPassword, err := util.HashPassword([]byte(req.NewPassword))
	if err != nil {
		return values.Error, "failed to hash new password", err
	}

	// store new password
	err = a.Deps.Repository.UpdateAdminPassword(ctx, executor.Id, hashPassword)
	if err != nil {
		return values.Error, "failed to update password", err
	}

	return values.Success, "updated password successfully", nil
}

func (a *API) UpdateAdminPassword(ctx context.Context, req model.UpdatePasswordReq) (string, string, error) {
	//encrypt new password
	hashPassword, err := util.HashPassword([]byte(req.Password))
	if err != nil {
		return values.Error, "failed to hash new password", err
	}
	err = a.Deps.Repository.UpdateAdminPassword(ctx, req.UserID, hashPassword)
	if err != nil {
		return values.Error, "failed to update password", err
	}

	return values.Success, "updated password successfully", nil
}

func (a *API) DoAdminRefreshToken(refreshToken string) (*model.RefreshTokenRes, string, string, error) {
	// get token details
	cClaims, err := util.ParseToken(refreshToken)
	if err != nil {
		return nil, values.NotAuthorised, "invalid refresh token", err
	}

	// get the session connected to refresh token

	AuthSession, err := a.Deps.Redis.GetRefreshSession(context.TODO(), cClaims.Role, cClaims.Subject, refreshToken)
	if err != nil {
		return nil, values.Error, "failed to get session", err
	}

	// delete that session
	err = a.Deps.Redis.DeleteAuthSession(context.TODO(), []string{fmt.Sprintf("ref-%s-%s-%s", cClaims.Role, cClaims.Subject, refreshToken)})
	if err != nil {
		return nil, values.Error, "failed to delete session", err
	}

	// check if sessionID is the same
	if cClaims.SessionId != AuthSession.SessionId {
		return nil, values.NotAuthorised, "invalid session Id", errors.New("user is logged in somewhere else")
	}

	// create new auth session
	_, token, status, message, err := a.CreateAuthToken(cClaims.Email, cClaims.Subject, cClaims.Role)
	if err != nil {
		return nil, status, message, err
	}

	refreshRes := &model.RefreshTokenRes{
		AccessToken:  token[0],
		RefreshToken: token[1],
	}

	return refreshRes, values.Success, "Authentication token reset successfully", nil
}
