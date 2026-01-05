package models

import "time"

type Comment struct {
	CommentID string    `json:"comment_id" db:"comment_id"`
	BodyText  string    `json:"body_text" db:"body_text"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	PostID    string    `json:"post_id" db:"post_id"`
	CreatedBy string    `json:"created_by" db:"created_by"`
}
