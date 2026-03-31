package routes

import (
	"collection-manager-backend/internal/handlers"
	"collection-manager-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", middleware.AuthMiddleware(), middleware.AdminOnly(), handlers.Register)
		authGroup.POST("/login", handlers.Login)
	}
}
