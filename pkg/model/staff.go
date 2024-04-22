package model

import (
	"github.com/google/uuid"
	"time"
)

type Staff struct {
	Id           uuid.UUID `json:"id,omitempty"`
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	Email        string    `json:"email,omitempty"`
	HashPassword string    `json:"-"`
	Status       string    `json:"status,omitempty"`
	State        string    `json:"state,omitempty"`
	Role         string    `json:"role,omitempty"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

type AuthStaff struct {
	Staff *Staff    `json:"staff,omitempty"`
	Auth  TokenInfo `json:"auth,omitempty"`
}

type ChangePasswordReq struct {
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}
type UpdatePasswordReq struct {
	Password string
	UserID   string // later get this from token
}

type RefreshTokenRes struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type PersistRes struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
