package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// creating user
func CreateUser(dbpool *pgxpool.Pool, username string, passwordHash string, email string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
			INSERT INTO users (username, password_hash, email) 
			VALUES ($1, $2, $3)
			RETURNING username, password_hash, user_id, created_at, email
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, username, passwordHash, email).Scan(
		&user.Username,
		&user.PasswordHash,
		&user.UserID,
		&user.CreatedAt,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Del User
func DeleteUser(dbpool *pgxpool.Pool, userID string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM users
		WHERE user_id = $1 
	`

	var commandTag, err = dbpool.Exec(ctx, query, userID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s not found", userID)
	}

	return nil
}

// Upd username -> Authentication to get userID is alrdy completed, password ws used to authenticate it
func UpdateUsername(dbpool *pgxpool.Pool, userID string, username string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE users
		SET username = $2
		WHERE user_id = $1 
		RETURNING user_id, username, password_hash, created_at, email
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, userID, username).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Upd password hash
// Take in userid + hashed password
// Clear reset token + token expiry
func UpdatePassword(dbpool *pgxpool.Pool, userID string, passwordHash string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE users
		SET password_hash = $2, password_reset_token_hash = NULL, password_reset_expire_at = NULL
		WHERE user_id = $1 
		RETURNING user_id, password_hash, username, created_at, email, password_reset_token_hash, password_reset_expire_at
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, userID, passwordHash).Scan(
		&user.UserID,
		&user.PasswordHash,
		&user.Username,
		&user.CreatedAt,
		&user.Email,
		&user.PasswordResetTokenHash,
		&user.PasswordResetExpireAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by username -> search
func GetUserByUsername(dbpool *pgxpool.Pool, username string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT username, user_id, password_hash, created_at, email
		FROM users
		WHERE username = $1
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, username).Scan(
		&user.Username,
		&user.UserID,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by userID -> fr chnging username
func GetUserByID(dbpool *pgxpool.Pool, userID string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT user_id, username, password_hash, created_at, email
		FROM users
		WHERE user_id = $1
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by email -> fr chngin password
func GetUserByEmail(dbpool *pgxpool.Pool, email string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT user_id, email, password_hash, created_at, username
		FROM users
		WHERE email = $1
	`
	var user models.User

	err := dbpool.QueryRow(ctx, query, email).Scan(
		&user.UserID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Username,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by token -> NEED CHECK
func GetUserByToken(dbpool *pgxpool.Pool, token string) (string, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT user_id
		FROM users
		WHERE password_reset_token_hash = $1
		 AND password_reset_expire_at > NOW()
	`
	var userID string
	err := dbpool.QueryRow(ctx, query, token).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func TokenStorage(dbpool *pgxpool.Pool, userID string, tokenHash string, expiresAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET password_reset_token_hash = $2,
		    password_reset_expire_at = $3
		WHERE user_id = $1
	`

	_, err := dbpool.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}
