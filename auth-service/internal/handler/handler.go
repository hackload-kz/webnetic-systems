package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"auth-svc/internal/domain/model"
	"auth-svc/internal/usecase"
)

type AuthHandler struct {
	service *usecase.AuthService
}

func NewAuthHandler(service *usecase.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// проверка пользователя 
func (h *AuthHandler) VerifyHandler(c *gin.Context) {
	var req model.VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	// TODO: Доделать
	ctx := context.TODO()
	response, err := h.service.VerifyUser(ctx ,req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	if !response.Success {
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// миддлварка для проверки jwt (for gateaway)
func (h *AuthHandler) ValidateTokenHandler(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Missing token",
		})
		return
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	claims := &model.AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.service.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user": gin.H{
			"id":         claims.UserID,
			"email":      claims.Email,
			"first_name": claims.FirstName,
			"last_name":  claims.LastName,
		},
	})
}
