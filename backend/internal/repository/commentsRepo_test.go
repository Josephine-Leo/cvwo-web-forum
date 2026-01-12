//KIV TESTING

package repository

import (
	"testing"

	"github.com/google/uuid"

	"context"
)

func TestCreateComment(t *testing.T) {
	userID := uuid.New().String()
	postID := uuid.New().String()

	// Create Test user
	_, err := testDB.Exec(context.Background(),
		"INSERT INTO users (user_id, username, password_hash) VALUES ($1, $2, $3)",
		userID, "testuser", "testhash",
	)
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}

	// Create Test post
	_, err = testDB.Exec(context.Background(),
		"INSERT INTO posts (post_id, title, created_by) VALUES ($1, $2, $3)",
		postID, "Test Post", userID,
	)
	if err != nil {
		t.Fatalf("failed to insert test post: %v", err)
	}

	// Create test comment
	comment, err := CreateComment(testDB, "hello world", "post-123", "user-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if comment.BodyText != "hello world" {
		t.Errorf("expected body_text 'hello world', got %s", comment.BodyText)
	}

	if comment.CommentID == "" {
		t.Errorf("expected comment_id to be set")
	}
}
