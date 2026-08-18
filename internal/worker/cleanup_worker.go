package worker

import (
	"context"
	"log"
	"time"

	"github.com/Golukpal/file-sharing/internal/service"
)

type CleanupWorker struct {
	service  *service.FileService
	interval time.Duration
}

func NewCleanupWorker(
	service *service.FileService,
	interval time.Duration,
) *CleanupWorker {

	return &CleanupWorker{
		service:  service,
		interval: interval,
	}
}

func (w *CleanupWorker) Start(ctx context.Context) {

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	log.Println("cleanup worker started")

	for {
		select {

		case <-ticker.C:

			log.Println("running cleanup...")

			if err := w.service.CleanupExpiredFiles(ctx); err != nil {
				log.Println(
					"cleanup failed:",
					err,
				)
			}

		case <-ctx.Done():

			log.Println(
				"cleanup worker stopped",
			)

			return
		}
	}
}