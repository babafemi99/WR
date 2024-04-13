package rest

import (
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"golang.org/x/crypto/bcrypt"
	"log"
)

func (a *API) DoSuperAdminLogin(req model.LoginReq) (*model.AdminAuthRes, string, string, error) {
	// validate req

	// fetch user with that email
	user, err := a.Deps.Repository.FindSuperAdminByEmail(req.Email)
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
	_, tokenArr, status, message, err := a.CreateAuthToken(user.Email, user.Id.String(), user.Role)
	if err != nil {
		return nil, status, message, err
	}

	// return response
	return &model.AdminAuthRes{
		Admin: &user,
		Auth: model.TokenInfo{
			Token:        tokenArr[0],
			RefreshToken: tokenArr[2],
		},
	}, values.Success, "log in successful", nil
}
