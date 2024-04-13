package redis

import (
	"time"
)

type AuthSession struct {
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
	Email                 string
	UserID                string
	SessionId             string
	// For of auth session - staff, admin or external entities
	For string
}

type ReferenceSession struct {
	RefreshToken string
	AccessToken  string
}
