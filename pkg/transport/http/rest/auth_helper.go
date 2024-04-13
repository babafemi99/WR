package rest

import (
	"context"
	"github.com/babafemi99/WR/internal/util"
	"github.com/babafemi99/WR/internal/values"
	"github.com/babafemi99/WR/pkg/services/redis"
	"github.com/oklog/ulid/v2"
	"strings"
	"time"
)

func (a *API) CreateAuthToken(email, id, authFor string) (redis.AuthSession, [3]string, string, string, error) {
	var (
		aT, rT, s      string
		err            error
		newAuthSession redis.AuthSession
	)

	sessionId := ulid.Make().String()

	switch strings.ToLower(authFor) {

	case "admin":
		aT, s, err = util.GenerateToken(id, authFor, email, sessionId, time.Now().Add(values.AccessTokenExpiry*time.Second))
		if err != nil {
			return redis.AuthSession{}, [3]string{}, values.Failed, "failed to generate token", err
		}

		rT, s, err = util.GenerateToken(id, authFor, email, sessionId, time.Now().Add(values.AccessTokenExpiry*time.Second))
		if err != nil {
			return redis.AuthSession{}, [3]string{}, values.Failed, "failed to generate token", err
		}

		newAuthSession = redis.AuthSession{
			Email:     email,
			SessionId: sessionId,
			For:       authFor,
			UserID:    id,
		}

		err = a.Deps.Redis.SetAuthSession(context.TODO(), newAuthSession, rT)
		if err != nil {
			return redis.AuthSession{}, [3]string{}, values.Failed, " failed to set session", err
		}

	case strings.ToLower("staff"):
		aT, s, err = util.GenerateToken(id, authFor, email, sessionId, time.Now().Add(values.StaffTokenExpiry*time.Hour))
		if err != nil {
			return redis.AuthSession{}, [3]string{}, values.Failed, "failed to generate token", err
		}

		newAuthSession = redis.AuthSession{
			Email:     email,
			SessionId: sessionId,
			For:       authFor,
			UserID:    id,
		}

		err = a.Deps.Redis.SetStaffSessionToken(context.TODO(), newAuthSession)
		if err != nil {
			return redis.AuthSession{}, [3]string{}, values.Failed, " failed to set session", err
		}

	case strings.ToLower("super"):
		aT, s, err = util.GenerateToken(id, authFor, email, sessionId, time.Now().Add(values.SuperTokenExpiry*time.Hour))
		if err != nil {
			return redis.AuthSession{}, [3]string{}, values.Failed, "failed to generate token", err
		}
		newAuthSession = redis.AuthSession{
			Email:     email,
			SessionId: sessionId,
			For:       authFor,
			UserID:    id,
		}

	default:
		return redis.AuthSession{}, [3]string{}, values.Failed, "invalid role", err

	}
	// admin

	// staff
	// - access token expiry
	// - no refresh sessions

	return newAuthSession, [3]string{aT, rT, s}, values.Success, "session created successfully", nil
}
