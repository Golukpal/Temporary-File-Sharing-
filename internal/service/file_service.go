package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/Golukpal/file-sharing/internal/repository"
)

type FileService struct {
	repo *repository.FileRepository
}

func NewFileService(repo *repository.FileRepository) *FileService {
	return &FileService{
		repo: repo,
	}
}

type UploadResult struct {
	ID           string    `json:"id"`
	OriginalName string    `json:"original_name"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (s *FileService) Upload(
	ctx context.Context,
	fileName string,
	fileSize int64,
	filePath string,
) (UploadResult, error) {

	id := uuid.New()

	storedName := id.String() + filepath.Ext(fileName)

	finalPath := filepath.Join("storage", storedName)

	err := os.Rename(filePath, finalPath)
	if err != nil {
		return UploadResult{}, fmt.Errorf("failed to store file: %w", err)
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	file := repository.File{
		ID:           id.String(),
		OriginalName: fileName,
		StoredName:   storedName,
		FilePath:     finalPath,
		FileSize:     fileSize,
		ExpiresAt:    expiresAt.Format(time.RFC3339),
	}

	if err := s.repo.Create(ctx, file); err != nil {
		_ = os.Remove(finalPath)

		return UploadResult{}, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return UploadResult{
		ID:           id.String(),
		OriginalName: fileName,
		ExpiresAt:    expiresAt,
	}, nil
}