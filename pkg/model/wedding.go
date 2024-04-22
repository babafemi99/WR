package model

import "time"

// Wedding Struct represents a wedding and all other related functionalities
type Wedding struct {
	Id          string    `json:"id,omitempty"`
	CoupleName  string    `json:"couple_name,omitempty"`
	Location    string    `json:"location,omitempty"`
	State       string    `json:"state"`
	Screen      int8      `json:"screen,omitempty"`
	Status      string    `json:"status,omitempty"`
	Link        string    `json:"link,omitempty"`
	WeddingId   string    `json:"wedding_id,omitempty"`
	WeddingDate time.Time `json:"wedding_date,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ModifiedAt  time.Time `json:"modified_at,omitempty"`
	ModifiedBy  string    `json:"modified_by,omitempty"`
}

type ToggleWeddingLink struct {
	Screen string `json:"screen,omitempty"`
}

type NewWeddingReq struct {
	CoupleName  string    `json:"couple_name,omitempty"`
	CoupleId    string    `json:"couple_id,omitempty"`
	State       string    `json:"state,omitempty"`
	WeddingId   string    `json:"wedding_id,omitempty"`
	Link        string    `json:"link,omitempty"`
	GuestLink   string    `json:"guest_link"`
	WeddingDate time.Time `json:"wedding_date,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type ToggleWeddingReq struct {
	TogglerId  string    `json:"toggler_id,omitempty"`
	WeddingId  string    `json:"wedding_id,omitempty"`
	Screen     int8      `json:"screen,omitempty"`
	Registry   string    `json:"registry,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

type Member struct {
	MemberEmail string `json:"member_email,omitempty"`
	MemberCode  string `json:"member_code,omitempty"`
	WeddingId   string `json:"wedding_id,omitempty"`
}

type NewWeddingRes struct {
	Link      string `json:"link,omitempty"`
	GuestLink string `json:"guest_link"`
	WeddingId string `json:"wedding_id,omitempty"`
}

type WeddingIdRes struct {
	CoupleName string `json:"couple_name,omitempty"`
	Location   string `json:"location,omitempty"`
	State      string `json:"state"`
	Screen     int8   `json:"screen,omitempty"`
	Status     string `json:"status,omitempty"`
}
