package models

import "time"

type Posts struct {
	PostID    string    `json:"post_id" db:"post_id"`
	Title     string    `json:"title" db:"title"`
	BodyText  string    `json:"body_text" db:"body_text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	TopicID   string    `json:"topic_id" db:"topic_id"`
	CreatedBy string    `json:"created_by" db:"created_by"`
}
