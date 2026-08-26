package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Golukpal/file-sharing/internal/config"
	"github.com/Golukpal/file-sharing/internal/database"
	"github.com/Golukpal/file-sharing/internal/handler"
	"github.com/Golukpal/file-sharing/internal/repository"
	"github.com/Golukpal/file-sharing/internal/service"
	"github.com/Golukpal/file-sharing/internal/worker"
)

func main() {
	cfg := config.Load()

	if err := os.MkdirAll("storage", 0755); err != nil {
		log.Fatal(
			"failed to create storage directory:",
			err,
		)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(
			"database connection failed:",
			err,
		)
	}

	defer db.Close()

	// Dependencies
	fileRepository := repository.NewFileRepository(db)

	fileService := service.NewFileService(
		fileRepository,
	)

	fileHandler := handler.NewFileHandler(
		fileService,
		cfg.MaxFileSize,
		cfg.AllowedFileTypes,
	)

	// Router
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			cfg.MaxFileSize+1024*1024,
		)

		c.Next()
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	router.POST(
		"/files",
		fileHandler.Upload,
	)

	router.GET(
		"/files/:id",
		fileHandler.Download,
	)

	// Worker context
	workerCtx, workerCancel := context.WithCancel(
		context.Background(),
	)

	defer workerCancel()

	// Cleanup worker
	cleanupWorker := worker.NewCleanupWorker(
		fileService,
		1*time.Minute,
	)

	go cleanupWorker.Start(workerCtx)

	// HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Printf(
			"server running on port %s",
			cfg.AppPort,
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatalf(
				"server failed: %v",
				err,
			)
		}
	}()

	// Wait for shutdown signal
	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-signalChan

	log.Println("shutdown signal received")

	// Stop accepting new requests
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)

	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf(
			"server shutdown failed: %v",
			err,
		)
	}

	// Stop background worker
	workerCancel()

	log.Println("server stopped")
}