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

	log.Printf(
		"cleanup worker started, interval=%s",
		w.interval,
	)

	for {
		select {
		case <-ticker.C:

			log.Println("cleanup started")

			if err := w.service.CleanupExpiredFiles(ctx); err != nil {
				log.Printf(
					"cleanup failed: %v",
					err,
				)
				continue
			}

			log.Println("cleanup completed")

		case <-ctx.Done():

			log.Println(
				"cleanup worker stopped",
			)

			return
		}
	}
}