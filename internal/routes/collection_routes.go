package routes

import (
	"collection-manager-backend/internal/handlers"
	"collection-manager-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterCollectionRoutes(r *gin.Engine) {
	collections := r.Group("/collections")
	collections.Use(middleware.AuthMiddleware())
	{
		collections.POST("", handlers.CreateCollection)
		collections.GET("", handlers.GetCollections)
		collections.PUT("/:id", handlers.UpdateCollection)
		collections.DELETE("/:id", handlers.DeleteCollection)
	}
}
