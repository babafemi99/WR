package model

import (
	"github.com/oklog/ulid/v2"
	"time"
)

type Staff struct {
	Id           ulid.ULID `json:"'id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	HashPassword string    `json:"hash_password"`
	Status       string    `json:"status"`
	State        string    `json:"state"`
	Role         string    `json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type AuthStaff struct {
	Staff *Staff    `json:"staff"`
	Auth  TokenInfo `json:"auth"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	UserID      string `json:"user_id"` // later get this from token
}
type UpdatePasswordReq struct {
	Password string
	UserID   string // later get this from token
}
