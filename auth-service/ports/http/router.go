package http

import (
	"github.com/gin-gonic/gin"

	"auth-svc/internal/handler"
)

func SetupRoutes(handler *handler.AuthHandler) *gin.Engine {
	r := gin.Default()
	
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/verify", handler.VerifyHandler)           // Для логина
			auth.POST("/validate", handler.ValidateTokenHandler) // Для gateway
		}
	}
	
	return r
}