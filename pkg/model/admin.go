package model

import (
	"github.com/google/uuid"
	"time"
)

// Admin represents an admin and their roles
type Admin struct {
	Id           uuid.UUID `json:"id,omitempty"`
	FirstName    string    `json:"first_name,omitempty"`
	LastName     string    `json:"last_name,omitempty"`
	Email        string    `json:"email,omitempty"`
	HashPassword string    `json:"-"`
	Role         string    `json:"role_id,omitempty"`
	Status       string    `json:"status,omitempty"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	DeletedAt    time.Time `json:"-"`
}

type BlockUserReq struct {
	Id        uuid.UUID `json:"id,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type LoginReq struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type TokenInfo struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type AdminAuthRes struct {
	Admin *Admin    `json:"admin,omitempty"`
	Auth  TokenInfo `json:"auth,omitempty"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Executor struct {
	Id    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role,omitempty"`
}
