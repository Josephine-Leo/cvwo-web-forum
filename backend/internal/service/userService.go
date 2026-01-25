package service

// What does my application allow a user to do?

import (
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/pkg/utils"

	//"backend/internal/pkg/utils"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(ctx context.Context, dbpool *pgxpool.Pool, username, email, Password string) (*models.User, error) {
	hashed, err := auth.HashPassword(Password)
	if err != nil {
		return nil, err
	}

	return repository.CreateUser(dbpool, username, hashed, email)
}

// use this fr ResetPasswordHandler, ResetPasswordHandler != UpdatePasswordHandler
func HandleResetPassword(ctx context.Context, dbpool *pgxpool.Pool, token string, newPassword string) (*models.User, error) {
	//hash provided token
	h := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(h[:])

	// Find userid by token hash
	userID, err := repository.GetUserByToken(dbpool, hashToken)
	if err != nil {
		return nil, fmt.Errorf("Invalid or expired reset token")
	}

	// hash new password
	hashed, err := auth.HashPassword(newPassword)
	if err != nil {
		return nil, err
	}

	return repository.UpdatePassword(dbpool, userID, hashed)
}

// Use password to upd -> use for ResetUsernameHandler func
func HandleForgotUsername(ctx context.Context, dbpool *pgxpool.Pool, email, password, newUsername string) (*models.User, error) {
	// get user by email
	user, err := repository.GetUserByEmail(dbpool, email)
	if err != nil {
		return nil, err
	}

	//compare hash password
	if !auth.ComparePassword(user.PasswordHash, password) {
		return nil, fmt.Errorf("Authenticated failed")
	}

	// upd + check for uniqueness of username
	updtedUser, err := repository.UpdateUsername(dbpool, user.UserID, newUsername)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("username already taken")
			}
		}
		return nil, err
	}

	return updtedUser, nil
}

func HandleForgotPassword(ctx context.Context, dbpool *pgxpool.Pool, email string) error {
	// Get User By Email
	user, err := repository.GetUserByEmail(dbpool, email)
	if err != nil {
		return err
	}

	// generate token
	tokenBytes := make([]byte, 32)
	_, err = rand.Read(tokenBytes)
	if err != nil {
		return err
	}
	token := hex.EncodeToString(tokenBytes)

	//hash provided token
	h := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(h[:])

	// Store hashed token + expiry -> Call repo
	err = repository.TokenStorage(dbpool, hashToken, user.UserID, time.Now().Add(10*time.Minute))
	if err != nil {
		return err
	}
	// Send email
	resetURL := fmt.Sprintf("http://localhost:3000/reset-password?token=%s", token)

	emailData := &utils.EmailData{
		URL:     resetURL,
		Subject: "Reset your password",
	}

	err = utils.SendEmail(user.Email, emailData, "reset_password.html")
	if err != nil {
		return err
	}

	return nil
}

// Handle Update Password?
