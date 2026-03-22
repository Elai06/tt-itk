package model

import "time"

type Wallet struct {
	UUID      string    `json:"uuid"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
