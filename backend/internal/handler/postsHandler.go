package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/repository"
	"net/http"
)

type CreatePostInput struct {
	Title    string `json:"title" binding:"required"`
	BodyText string `json:"body_text" binding:"required"`
}
type UpdatePostInput struct {
	BodyText *string `json:"body_text"`
}

// Create Post Handler
func CreatePostHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		//Authenticate users
		userIDInterface, exists := c.Get("user_id") // this frm bwt authentication NT repository

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		// Get TopicID
		topicID := c.Param("id")
		if topicID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "topic_id is required"})
			return
		}

		//creatin request
		var input CreatePostInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// calling repository
		post, err := repository.CreatePost(dbpool, input.BodyText, input.Title, userID, topicID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, post)
	}
}

// Update Post Handler
func UpdatePostHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)
		postID := c.Param("id")
		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Post ID"})
			return
		}

		var input UpdatePostInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.BodyText == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Update cannot be empty"})
			return
		}

		existing, err := repository.GetPostByID(dbpool, postID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Post Not Found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		bodyText := existing.BodyText
		if input.BodyText != nil {
			bodyText = *input.BodyText
		}

		comment, err := repository.UpdatePost(dbpool, postID, bodyText, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, comment)
	}
}

// Delete Post Handler -> the logic wld be authentication?

func DeletePostHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		postID := c.Param("id")

		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		err := repository.DeletePost(dbpool, postID, userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	}
}

// Get all post by user ID
func GetPostByIDHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		posts, err := repository.GetPostByID(dbpool, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, posts)
	}
}

// Get all post by topic ID
func GetPostByTopicHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID := c.Param("topic_id")
		if topicID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "topic_id is required"})
			return
		}

		posts, err := repository.GetPostByTopic(dbpool, topicID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, posts)
	}
}

// Get single post
func GetPostByPostIDHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("post_id")
		post, err := repository.GetPostByPostID(dbpool, postID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, post)
	}
}
