package model

import (
	"time"
)

type User struct {
	ID        string    `json:"ID"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}
