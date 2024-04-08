package model

import "time"

// Wedding Struct represents a wedding and all other related functionalities
type Wedding struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	Screen     string    `json:"screen"`
	Live       bool      `json:"live"`
	Link       string    `json:"link"`
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
