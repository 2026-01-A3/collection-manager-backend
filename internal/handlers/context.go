package handlers

import (
	"net/http"

	"collection-manager-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// actorFromContext lê user_id e user_role setados pelo AuthMiddleware.
// Retorna (userID, isAdmin, ok). Se ok=false, já respondeu 401.
func actorFromContext(c *gin.Context) (uint, bool, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return 0, false, false
	}
	id, ok := v.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário inválido no contexto"})
		return 0, false, false
	}

	role, _ := c.Get("user_role")
	isAdmin := role == models.AdminRole

	return id, isAdmin, true
}
