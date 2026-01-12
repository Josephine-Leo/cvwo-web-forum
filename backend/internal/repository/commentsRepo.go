package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Funcs
// CreateComment, DeleteComment, UpdateComment, GetCommentByID, GetCommentByPost

// create comment
func CreateComment(dbpool *pgxpool.Pool, BodyText string, PostID string, CreatedBy string) (*models.Comment, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
			INSERT INTO comments (body_text, post_id, created_by) 
			VALUES ($1, $2, $3)
			RETURNING body_text, post_id, created_by, created_at, updated_at, comment_id 
	`

	var comment models.Comment

	var err error = dbpool.QueryRow(ctx, query, BodyText, PostID, CreatedBy).Scan(
		&comment.BodyText,
		&comment.PostID,
		&comment.CreatedBy,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.CommentID,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// Del comments
func DeleteComment(dbpool *pgxpool.Pool, CommentID string, CreatedBy string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM comments
		WHERE comment_id = $1 AND created_by = $2 
	`

	var commandTag, err = dbpool.Exec(ctx, query, CommentID, CreatedBy)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("post with id %s not found", CommentID)
	}

	return nil
}

// Upd comments
func UpdateComment(dbpool *pgxpool.Pool, commentID string, bodyText string, createdBy string) (*models.Comment, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE comments
		SET body_text = $2, updated_at = CURRENT_TIMESTAMP
		WHERE comment_id = $1 AND created_by = $3
		RETURNING comment_id, body_text, created_at, updated_at, post_id, created_by
	`

	var comment models.Comment

	var err error = dbpool.QueryRow(ctx, query, commentID, bodyText, createdBy).Scan(
		&comment.CommentID,
		&comment.BodyText,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.PostID,
		&comment.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// Get comments by UserID -> See what comments u hv
func GetCommentByID(dbpool *pgxpool.Pool, createdBy string) (*models.Comment, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT created_by, comment_id, body_text, created_at, updated_at, post_id
		FROM comments
		WHERE created_by = $1
	`

	var comment models.Comment

	var err error = dbpool.QueryRow(ctx, query, createdBy).Scan(
		&comment.CreatedBy,
		&comment.CommentID,
		&comment.BodyText,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.PostID,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}

// Get comments by PostID -> See what comments a post hv
func GetCommentByPost(dbpool *pgxpool.Pool, postID string) (*models.Comment, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT post_id, created_by, comment_id, body_text, created_at, updated_at
		FROM comments
		WHERE post_id = $1
	`

	var comment models.Comment

	var err error = dbpool.QueryRow(ctx, query, postID).Scan(
		&comment.PostID,
		&comment.CreatedBy,
		&comment.CommentID,
		&comment.BodyText,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &comment, nil
}
