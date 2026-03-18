package main

import (
	"log"
	"time"

	"collection-manager-backend/internal/database"
	"collection-manager-backend/internal/routes"
	"collection-manager-backend/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("aviso: arquivo .env não carregado, usando variáveis de ambiente do sistema")
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	if err := storage.InitCategoryStorage(db); err != nil {
		log.Fatalf("erro ao inicializar storage de categorias: %v", err)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterCategoryRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
