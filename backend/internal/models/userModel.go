package models

import (
	"time"
)

type User struct {
	UserID                 string    `json:"user_id" db:"user_id"`
	Username               string    `json:"username" db:"username"`
	PasswordHash           string    `json:"password_hash" db:"password_hash"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	Email                  string    `json:"email" db:"email"`
	PasswordResetTokenHash string    `json:"password_reset_token_hash" db:"password_reset_token_hash"`
	PasswordResetExpireAt  time.Time `json:"password_reset_expire_at" db:"Password_reset_expire_at"`
}
