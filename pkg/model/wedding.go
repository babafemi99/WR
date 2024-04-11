package model

import "time"

// Wedding Struct represents a wedding and all other related functionalities
type Wedding struct {
	Id         string    `json:"id"`
	CoupleName string    `json:"couple_name"`
	Location   string    `json:"location"`
	Screen     int8      `json:"screen"`
	Live       bool      `json:"live"`
	Link       string    `json:"link"`
	WeddingId  string    `json:"wedding_id"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	ModifiedBy string    `json:"modified_by"`
}

func (w Wedding) GetLink() string {
	return w.Link
}

func (w Wedding) IsLive() bool {
	return w.Live
}

type ToggleWeddingLink struct {
	Screen string `json:"screen"`
}

type NewWeddingReq struct {
	CoupleName string    `json:"couple_name"`
	CoupleId   string    `json:"couple_id"`
	Location   string    `json:"location"`
	WeddingId  string    `json:"wedding_id"`
	Link       string    `json:"link"`
	CreatedAt  time.Time `json:"created_at"`
}

type ToggleWeddingReq struct {
	TogglerId  string    `json:"toggler_id"`
	WeddingId  string    `json:"wedding_id"`
	Screen     int8      `json:"screen"`
	Registry   string    `json:"registry"`
	ModifiedAt time.Time `json:"modified_at"`
}

type Member struct {
	MemberEmail string `json:"member_email"`
	MemberCode  string `json:"member_code"`
	WeddingId   string `json:"wedding_id"`
}
