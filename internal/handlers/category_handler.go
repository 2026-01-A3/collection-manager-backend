package handlers

import (
	"fmt"
	"net/http"

	"collection-manager-backend/internal/storage"

	"github.com/gin-gonic/gin"
)

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func CreateCategory(c *gin.Context) {
	var input CreateCategoryInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	category, err := storage.AddCategory(c.Request.Context(), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao salvar categoria",
		})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func GetCategories(c *gin.Context) {
	categories, err := storage.GetCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao buscar categorias",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}

type UpdateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id := 0
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	var input UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	category, err := storage.UpdateCategory(c.Request.Context(), id, input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao atualizar categoria",
		})
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id := 0
	if _, err := fmt.Sscanf(idParam, "%d", &id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	err := storage.DeleteCategory(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erro ao deletar categoria",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Categoria deletada com sucesso",
	})
}
