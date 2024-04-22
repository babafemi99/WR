package rest

import (
	"context"
	"errors"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/model"
	"log"
	"net/http"
	"strings"
)

// AuthenticateStaff  checks if staff token is valid
func (a *API) AuthenticateStaff(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authorization) != 2 {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, "not-authorized")
			return
		}

		token, err := util.ValidateToken(authorization[1])
		if err != nil {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, err.Error())
			return
		}

		// check kind of error if it is error no row or system error
		session, err := a.Deps.Redis.GetAuthSession(context.TODO(), token.Role, token.Subject)
		if err != nil {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, err.Error())
			return
		}

		// check if sessionID is the same
		if token.SessionId != session.SessionId {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, "you are logged in somewhere else")
			return
		}

		if strings.ToLower(token.Role) != "staff" {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, " who you ?? you no supposed dey here")
			return
		}

		executor := model.Executor{
			Id:    token.Subject,
			Email: token.Email,
			Role:  token.Role,
		}

		// add executor
		ctx := context.WithValue(r.Context(), values.Executor, executor)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// AuthenticateAdmin checks if admin token is valid
func (a *API) AuthenticateAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authorization) != 2 {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, "not-authorized")
			return
		}

		token, err := util.ValidateToken(authorization[1])
		if err != nil {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, err.Error())
			return
		}

		if token.Role == "admin" {
			_, err = a.Deps.Redis.GetAuthSession(r.Context(), token.Role, token.Subject)
			if err != nil {
				writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, err.Error())
				return
			}
		}

		if strings.ToLower(token.Role) != "admin" && strings.ToLower(token.Role) != "super" {
			log.Println(token.Role)
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, " who you ?? you no supposed dey here")
			return
		}

		executor := model.Executor{
			Id:    token.Subject,
			Email: token.Email,
			Role:  token.Role,
		}

		// add executor
		ctx := context.WithValue(r.Context(), values.Executor, executor)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// AuthenticateSuperAdmin  checks if boss token is valid
func AuthenticateSuperAdmin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.Split(r.Header.Get("Authorization"), " ")
		if len(authorization) != 2 {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, "invalid token format")
			return
		}

		token, err := util.ValidateToken(authorization[1])
		if err != nil {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, err.Error())
			return
		}

		if strings.ToLower(token.Role) != "super" {
			writeErrorResponse(w, errors.New(values.NotAuthorised), values.NotAuthorised, " who you ?? you no supposed dey here")
			return
		}

		executor := model.Executor{
			Id:    token.Subject,
			Email: token.Email,
			Role:  token.Role,
		}

		// add executor
		ctx := context.WithValue(r.Context(), values.Executor, executor)

		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
