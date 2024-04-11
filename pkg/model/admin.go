package model

import (
	"github.com/oklog/ulid/v2"
	"time"
)

// Admin represents an admin and their roles
type Admin struct {
	Id           ulid.ULID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	HashPassword string    `json:"-"`
	Role         string    `json:"role_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

type BlockUserReq struct {
	Id        ulid.ULID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenInfo struct {
	Token           string `json:"token"`
	TokenExpiryTime time.Time

	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`
}

type AdminAuthRes struct {
	Admin *Admin    `json:"admin"`
	Auth  TokenInfo `json:"auth"`
}
