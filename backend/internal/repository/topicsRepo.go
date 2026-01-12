package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Funcs
// CreateTopic, DeleteTopic, GetTopicByUserID
// GetTopicByTopicID, GetAllTopics

// create topic -> Return order nd to match scan order
func CreateTopic(dbpool *pgxpool.Pool, Title string, Subtitle string, CreatedBy string) (*models.Topic, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
			INSERT INTO topics (title, subtitle, created_by) 
			VALUES ($1, $2, $3)
			RETURNING title, subtitle, created_by, created_at, topic_id
	`

	var topic models.Topic

	var err error = dbpool.QueryRow(ctx, query, Title, Subtitle, CreatedBy).Scan(
		&topic.Title,
		&topic.Subtitle,
		&topic.CreatedBy,
		&topic.CreatedAt,
		&topic.TopicID,
	)

	if err != nil {
		return nil, err
	}

	return &topic, nil
}

// Del topic
func DeleteTopic(dbpool *pgxpool.Pool, TopicID string, CreatedBy string) error {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		DELETE FROM topics
		WHERE topic_id = $1 AND created_by = $2 
	`

	var commandTag, err = dbpool.Exec(ctx, query, TopicID, CreatedBy)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("topic with id %s not found", TopicID)
	}

	return nil
}

// Get topic by user -> see what topic each user created
func GetTopicByUserID(dbpool *pgxpool.Pool, createdBy string) (*models.Topic, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT created_by, topic_id, title, subtitle, created_at
		FROM topics
		WHERE created_by = $1
	`

	var topic models.Topic

	var err error = dbpool.QueryRow(ctx, query, createdBy).Scan(
		&topic.CreatedBy,
		&topic.TopicID,
		&topic.Title,
		&topic.Subtitle,
		&topic.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &topic, nil
}

// Get topic by ID -> to get specific topic
func GetTopicByTopicID(dbpool *pgxpool.Pool, topicID string) (*models.Topic, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT topic_id, title, subtitle, created_at, created_by
		FROM topics
		WHERE topic_id = $1
	`

	var topic models.Topic

	var err error = dbpool.QueryRow(ctx, query, topicID).Scan(
		&topic.TopicID,
		&topic.Title,
		&topic.Subtitle,
		&topic.CreatedAt,
		&topic.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &topic, nil
}

// Get all topics -> Show on frontend
func GetAllTopics(dbpool *pgxpool.Pool) ([]models.Topic, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query string = `
		SELECT topic_id, title, subtitle, created_at, created_by
		FROM topics
		ORDER BY created_at DESC
	`

	var rows, err = dbpool.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var topics []models.Topic = []models.Topic{}

	for rows.Next() {
		var topic models.Topic

		err = rows.Scan(
			&topic.TopicID,
			&topic.Title,
			&topic.Subtitle,
			&topic.CreatedAt,
			&topic.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		topics = append(topics, topic)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return topics, nil
}
