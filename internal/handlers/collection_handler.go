package handlers

import (
	"fmt"
	"net/http"

	"collection-manager-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type CreateCollectionInput struct {
	Name       string `json:"name" binding:"required"`
	CategoryID int    `json:"category_id" binding:"required"`
	ImageURL   string `json:"image_url"`
}

func CreateCollection(c *gin.Context) {
	var input CreateCollectionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	collection, err := storage.AddCollection(c.Request.Context(), input.Name, input.CategoryID, input.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao salvar coleção",
		})
		return
	}

	c.JSON(http.StatusCreated, collection)
}

func GetCollections(c *gin.Context) {
	collections, err := storage.GetCollections(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar coleções",
		})
		return
	}

	c.JSON(http.StatusOK, collections)
}

type UpdateCollectionInput struct {
	Name       string `json:"name" binding:"required"`
	CategoryID int    `json:"category_id" binding:"required"`
	ImageURL   string `json:"image_url"`
}

func UpdateCollection(c *gin.Context) {
	idParam := c.Param("id")
	id := 0
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	var input UpdateCollectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	collection, err := storage.UpdateCollection(c.Request.Context(), id, input.Name, input.CategoryID, input.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar coleção",
		})
		return
	}

	c.JSON(http.StatusOK, collection)
}

func DeleteCollection(c *gin.Context) {
	idParam := c.Param("id")
	id := 0
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	err := storage.DeleteCollection(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao deletar coleção",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Coleção deletada com sucesso",
	})
}
