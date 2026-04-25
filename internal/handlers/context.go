package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func userIDFromContext(c *gin.Context) (uint, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return 0, false
	}
	id, ok := v.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário inválido no contexto"})
		return 0, false
	}
	return id, true
}
