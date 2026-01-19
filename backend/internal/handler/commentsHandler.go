package handler

import (
	"backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// What users get to input
type CreateCommentInput struct {
	BodyText string `json:"body_text" binding:"required"`
}

type UpdateCommentInput struct {
	BodyText *string `json:"body_text"`
}

// CreateCommentHandler, UpdateCommentHandler, DeleteCommentHandler
func CreateCommentHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		//Authenticate users
		userIDInterface, exists := c.Get("user_id") // this frm bwt authentication NT repository

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		// Get PostID
		postID := c.Param("id")
		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
			return
		}

		//creatin request
		var input CreateCommentInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// calling repository
		comment, err := repository.CreateComment(dbpool, input.BodyText, postID, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, comment)
	}
}

func UpdateCommentHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)
		commentID := c.Param("id")
		if commentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Comment ID"})
			return
		}

		var input UpdateCommentInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if input.BodyText == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Update cannot be empty"})
			return
		}

		existing, err := repository.GetCommentByID(dbpool, commentID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Comment Not Found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		bodyText := existing.BodyText
		if input.BodyText != nil {
			bodyText = *input.BodyText
		}

		comment, err := repository.UpdateComment(dbpool, commentID, bodyText, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, comment)
	}
}

func DeleteCommentHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		commentID := c.Param("id")

		if commentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		err := repository.DeleteComment(dbpool, commentID, userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
	}
}

// Get all comment by user ID
func GetCommentByIDHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		comments, err := repository.GetCommentByID(dbpool, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, comments)
	}
}

// Get all comments by post ID
func GetCommentByPostHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID := c.Param("post_id")
		if postID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "post_id is required"})
			return
		}

		comments, err := repository.GetCommentByPost(dbpool, postID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, comments)
	}
}

// idk if nd to get single comment
