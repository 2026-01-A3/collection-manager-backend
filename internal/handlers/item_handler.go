package handlers

import (
	"errors"
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
	userID, isAdmin, ok := actorFromContext(c)
	if !ok {
		return
	}

	var input CreateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	item, err := storage.AddItem(c.Request.Context(), userID, isAdmin, input.Name, input.Price, input.CollectionID, toPayload(input.BinaryObject))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coleção não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar item"})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func GetItemsByCollection(c *gin.Context) {
	userID, isAdmin, ok := actorFromContext(c)
	if !ok {
		return
	}

	idParam := c.Query("collection_id")
	collectionID := 0
	if _, err := fmt.Sscanf(idParam, "%d", &collectionID); err != nil || collectionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "collection_id inválido"})
		return
	}

	items, err := storage.GetItemsByCollection(c.Request.Context(), userID, isAdmin, collectionID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coleção não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar itens"})
		return
	}

	c.JSON(http.StatusOK, items)
}

func UpdateItem(c *gin.Context) {
	userID, isAdmin, ok := actorFromContext(c)
	if !ok {
		return
	}

	idParam := c.Param("id")
	id := 0
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var input UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	item, err := storage.UpdateItem(c.Request.Context(), userID, isAdmin, id, input.Name, input.Price, input.CollectionID, toPayload(input.BinaryObject))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item ou coleção não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func DeleteItem(c *gin.Context) {
	userID, isAdmin, ok := actorFromContext(c)
	if !ok {
		return
	}

	idParam := c.Param("id")
	id := 0
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := storage.DeleteItem(c.Request.Context(), userID, isAdmin, id); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Item não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao deletar item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item deletado com sucesso"})
}
