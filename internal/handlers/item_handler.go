package handlers

import (
	"fmt"
	"net/http"

	"collection-manager-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type CreateItemInput struct {
	Name         string             `json:"name" binding:"required"`
	Price        float64            `json:"price"`
	CollectionID int                `json:"collection_id" binding:"required"`
	BinaryObject *BinaryObjectInput `json:"binary_object"`
}

type UpdateItemInput struct {
	Name         string             `json:"name" binding:"required"`
	Price        float64            `json:"price"`
	CollectionID int                `json:"collection_id" binding:"required"`
	BinaryObject *BinaryObjectInput `json:"binary_object"`
}

func CreateItem(c *gin.Context) {
	var input CreateItemInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	item, err := storage.AddItem(c.Request.Context(), input.Name, input.Price, input.CollectionID, toPayload(input.BinaryObject))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao salvar item",
		})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func GetItemsByCollection(c *gin.Context) {
	idParam := c.Query("collection_id")
	collectionID := 0
	if _, err := fmt.Sscanf(idParam, "%d", &collectionID); err != nil || collectionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "collection_id inválido",
		})
		return
	}

	items, err := storage.GetItemsByCollection(c.Request.Context(), collectionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar itens",
		})
		return
	}

	c.JSON(http.StatusOK, items)
}

func UpdateItem(c *gin.Context) {
	idParam := c.Param("id")
	id := 0
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	var input UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	item, err := storage.UpdateItem(c.Request.Context(), id, input.Name, input.Price, input.CollectionID, toPayload(input.BinaryObject))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar item",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func DeleteItem(c *gin.Context) {
	idParam := c.Param("id")
	id := 0
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	if err := storage.DeleteItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao deletar item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item deletado com sucesso",
	})
}
