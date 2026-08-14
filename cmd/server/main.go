package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/Golukpal/file-sharing/internal/config"
	"github.com/Golukpal/file-sharing/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal("database connection failed:", err)
	}

	defer db.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	log.Println("server running on port", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}