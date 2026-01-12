package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Funcs
// CreateUser, DeleteUser, UpdateUsername, UpdatePassword
// GetUserByUsername, GetUserByID

// creating user EDIT TO MATCH SEQ IN MODEL
func CreateUser(dbpool *pgxpool.Pool, Username string, PasswordHash string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
			INSERT INTO users (username, password_hash) 
			VALUES ($1, $2)
			RETURNING username, password_hash, user_id, created_at
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, Username, PasswordHash).Scan(
		&user.Username,
		&user.PasswordHash,
		&user.UserID,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Del User [KIV]
func DeleteUser(dbpool *pgxpool.Pool, UserID string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM users
		WHERE user_id = $1 
	`

	var commandTag, err = dbpool.Exec(ctx, query, UserID)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id %s not found", UserID)
	}

	return nil
}

// Upd username
func UpdateUsername(dbpool *pgxpool.Pool, userID string, username string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE users
		SET username = $2
		WHERE user_id = $1 
		RETURNING user_id, username, password_hash, created_at
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, userID, username).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Upd password hash -> whn user chnge password
func UpdatePassword(dbpool *pgxpool.Pool, userID string, passwordHash string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE users
		SET password_hash = $2
		WHERE user_id = $1 
		RETURNING user_id, password_hash, username, created_at
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, userID, passwordHash).Scan(
		&user.UserID,
		&user.PasswordHash,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Get user by username -> Log in
func GetUserByUsername(dbpool *pgxpool.Pool, username string) (*models.User, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT username, user_id, password_hash, created_at
		FROM users
		WHERE username = $1
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, username).Scan(
		&user.Username,
		&user.UserID,
		&user.PasswordHash,
		&user.CreatedAt,
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
		SELECT user_id, username, password_hash, created_at
		FROM users
		WHERE user_id = $1
	`

	var user models.User

	var err error = dbpool.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
