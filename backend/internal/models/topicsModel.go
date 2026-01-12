package models

import (
	"time"
)

type Topic struct {
	TopicID   string    `json:"topic_id" db:"topic_id"`
	CreatedBy string    `json:"created_by" db:"created_by"`
	Title     string    `json:"title" db:"title"`
	Subtitle  string    `json:"subtitle" db:"subtitle"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
