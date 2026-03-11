package routes

import (
	"collection-manager-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(router *gin.Engine) {
	categoryGroup := router.Group("/categories")
	{
		categoryGroup.POST("", handlers.CreateCategory)
		categoryGroup.GET("", handlers.GetCategories)
	}
}
