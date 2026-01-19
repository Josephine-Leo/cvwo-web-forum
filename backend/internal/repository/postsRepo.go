package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Funcs
// CreatePost, DeletePost, UpdatePost, GetPostByID, GetPostByTopic

// creating user accept all 4
func CreatePost(dbpool *pgxpool.Pool, TopicID string, Title string, BodyText string, CreatedBy string) (*models.Post, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
			INSERT INTO posts (topic_id, title, body_text, created_by) 
			VALUES ($1, $2, $3, $4)
			RETURNING topic_id, title, body_text, created_at, created_by, updated_at, post_id
	`

	var post models.Post

	var err error = dbpool.QueryRow(ctx, query, TopicID, Title, CreatedBy, BodyText).Scan(
		&post.TopicID,
		&post.Title,
		&post.BodyText,
		&post.CreatedAt,
		&post.CreatedBy,
		&post.UpdatedAt,
		&post.PostID,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Del post
func DeletePost(dbpool *pgxpool.Pool, PostID string, CreatedBy string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM posts
		WHERE post_id = $1 AND created_by = $2 
	`

	var commandTag, err = dbpool.Exec(ctx, query, PostID, CreatedBy)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("post with id %s not found", PostID)
	}

	return nil
}

// Upd Post body text (cnt upd title)
func UpdatePost(dbpool *pgxpool.Pool, postID string, bodyText string, createdBy string) (*models.Post, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		UPDATE posts
		SET body_text = $2, updated_at = CURRENT_TIMESTAMP
		WHERE post_id = $1 AND created_by = $3
		RETURNING post_id, body_text, created_by, updated_at, created_at, title, topic_id
	`

	var post models.Post

	var err error = dbpool.QueryRow(ctx, query, postID, bodyText, createdBy).Scan(
		&post.PostID,
		&post.BodyText,
		&post.CreatedBy,
		&post.UpdatedAt,
		&post.CreatedAt,
		&post.Title,
		&post.TopicID,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Get ALL post by UserID -> See what you posted
func GetPostByID(dbpool *pgxpool.Pool, createdBy string) (*models.Post, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT post_id, title, body_text, created_at, updated_at, topic_id, created_by
		FROM posts
		WHERE created_by = $1
	`

	var post models.Post

	var err error = dbpool.QueryRow(ctx, query, createdBy).Scan(
		&post.PostID,
		&post.Title,
		&post.BodyText,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.TopicID,
		&post.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Get ALL post by Topic -> See what posts in a topic
func GetPostByTopic(dbpool *pgxpool.Pool, topicID string) (*models.Post, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT topic_id, post_id, title, body_text, created_at, updated_at, created_by
		FROM posts
		WHERE topic_id = $1
	`

	var post models.Post

	var err error = dbpool.QueryRow(ctx, query, topicID).Scan(
		&post.TopicID,
		&post.PostID,
		&post.Title,
		&post.BodyText,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

// Get single post
func GetPostByPostID(dbpool *pgxpool.Pool, postID string) (*models.Post, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT post_id, title, body_text, created_at, updated_at, topic_id, created_by
		FROM posts
		WHERE post_id = $1
	`

	var post models.Post

	var err error = dbpool.QueryRow(ctx, query, post).Scan(
		&post.PostID,
		&post.Title,
		&post.BodyText,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.TopicID,
		&post.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &post, nil
}
