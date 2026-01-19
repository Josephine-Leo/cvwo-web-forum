package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/repository"
	"net/http"
)

type CreateTopicInput struct {
	Title    string `json:"title" binding:"required"`
	Subtitle string `json:"subtitle"`
}

// create topic
func CreateTopicHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		//Authenticate users
		userIDInterface, exists := c.Get("user_id") // this frm bwt authentication NT repository

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		//creatin request
		var input CreateTopicInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// calling repository
		post, err := repository.CreateTopic(dbpool, input.Title, input.Subtitle, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, post)
	}
}

// delete topic -> User dont input anyth
func DeleteTopicHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		topicID := c.Param("id")

		if topicID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
			return
		}

		err := repository.DeleteTopic(dbpool, topicID, userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Topic deleted successfully"})
	}
}

// Get all topics
func GetAllTopicsHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		topics, err := repository.GetAllTopics(dbpool)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, topics)
	}
}

// Get all topics by user
func GetTopicByUserIDHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		topics, err := repository.GetTopicByUserID(dbpool, userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, topics)
	}
}

// Get single topic
func GetTopicByTopicIDHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID := c.Param("topic_id")
		topic, err := repository.GetTopicByTopicID(dbpool, topicID)

		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Topic not found"})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, topic)
	}
}
