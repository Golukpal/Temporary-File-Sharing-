package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Golukpal/file-sharing/internal/config"
	"github.com/Golukpal/file-sharing/internal/database"
	"github.com/Golukpal/file-sharing/internal/handler"
	"github.com/Golukpal/file-sharing/internal/repository"
	"github.com/Golukpal/file-sharing/internal/service"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll("storage", 0755); err != nil {
		log.Fatal("failed to create storage directory:", err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}

	defer db.Close()

	fileRepository := repository.NewFileRepository(db)

	fileService := service.NewFileService(fileRepository)

	fileHandler := handler.NewFileHandler(fileService)

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.POST("/files", fileHandler.Upload)

	log.Println("server running on port", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}