package routes

import (
	"collection-manager-backend/internal/handlers"
	"collection-manager-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(router *gin.Engine) {
	categoryGroup := router.Group("/categories")
	categoryGroup.Use(middleware.AuthMiddleware())
	{
		categoryGroup.POST("", handlers.CreateCategory)
		categoryGroup.GET("", handlers.GetCategories)
	}
}
