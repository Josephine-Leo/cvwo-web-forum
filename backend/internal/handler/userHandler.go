package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/repository"
	"backend/internal/service"
	"net/http"
	"strings"
)

// NOTE
// Need to upd Handler to include service layer in password n username funcs
type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

// Authenticated Users
type UpdateUsernameInput struct {
	Username string `json:"username"`
}

type UpdatePasswordInput struct {
	NewPassword string `json:"new_password"`
}

// Unauthenticated Users
type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type ForgotUsernameInput struct {
	Password    string `json:"password"`
	Email       string `json:"email"`
	NewUsername string `json:"new_username"`
}

// create user Did wrongly i think
func CreateUserHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input CreateUserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		// service layer
		user, err := service.CreateUser(c.Request.Context(), dbpool, input.Username, input.Email, input.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":       user.UserID,
			"username": user.Username,
			"email":    user.Email,
		})
	}
}

// del user -> No input
func DeleteUserHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")

		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id not found in context"})
			return
		}

		userID := userIDInterface.(string)

		err := repository.DeleteUser(dbpool, userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	}
}

// upd username -> This user is authenticated alrdy
func UpdateUsernameHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
			return
		}

		userID, ok := userIDValue.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
			return
		}

		var input UpdateUsernameInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if strings.TrimSpace(input.Username) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username cannot be empty"})
			return
		}

		user, err := repository.UpdateUsername(dbpool, userID, input.Username)
		if err != nil {
			if err == pgx.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":       user.UserID,
			"username": user.Username,
		})
	}
}

// forget password -> use email to get userID (unauthenticated)
func ForgotPasswordHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ForgotPasswordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if strings.TrimSpace(input.Email) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email cannot be empty"})
			return
		}

		err := service.HandleForgotPassword(c.Request.Context(), dbpool, input.Email)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "email verification sent",
		})
	}
}

// forget username -> use password + email to get userID (unauthenticated)
func ForgotUsernameHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ForgotUsernameInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if strings.TrimSpace(input.Email) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email cannot be empty"})
			return
		}

		if strings.TrimSpace(input.NewUsername) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "new username cannot be empty"})
			return
		}

		if strings.TrimSpace(input.Password) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password cannot be empty"})
			return
		}

		_, err := service.HandleForgotUsername(c.Request.Context(), dbpool, input.Email, input.NewUsername, input.Password)

		if err != nil {
			switch err.Error() {
			case "authentication failed":
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			case "username already taken":
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "username updated successfully",
		})

	}
}

// Reset Password Handler (aftr authenticated by forgot password)
func ResetPasswordHandler(dbpool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ResetPasswordInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if strings.TrimSpace(input.Token) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "reset token is required"})
			return
		}

		if strings.TrimSpace(input.NewPassword) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "password cannot be empty"})
			return
		}

		_, err := service.HandleResetPassword(c.Request.Context(), dbpool, input.NewPassword, input.Token)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "password reset is successful",
		})
	}
}

// https://codevoweb.com/api-golang-gin-gonic-mongodb-forget-reset-password/
