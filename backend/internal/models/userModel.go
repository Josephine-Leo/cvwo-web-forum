package models

import (
	"time"
)

type User struct {
	UserID       string    `json:"user_id" db:"user_id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
