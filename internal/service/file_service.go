package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/Golukpal/file-sharing/internal/repository"
)

var ErrFileNotFound = errors.New("file not found")
var ErrFileExpired = errors.New("file has expired")

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
		return UploadResult{}, fmt.Errorf(
			"failed to store file: %w",
			err,
		)
	}

	expiresAt := time.Now().Add(24 * time.Hour)

	file := repository.File{
		ID:           id.String(),
		OriginalName: fileName,
		StoredName:   storedName,
		FilePath:     finalPath,
		FileSize:     fileSize,
		ExpiresAt:    expiresAt,
	}

	if err := s.repo.Create(ctx, file); err != nil {
		_ = os.Remove(finalPath)

		return UploadResult{}, fmt.Errorf(
			"failed to save file metadata: %w",
			err,
		)
	}

	return UploadResult{
		ID:           id.String(),
		OriginalName: fileName,
		ExpiresAt:    expiresAt,
	}, nil
}

func (s *FileService) GetFile(
	ctx context.Context,
	id string,
) (repository.File, error) {

	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrFileNotFound) {
			return repository.File{}, ErrFileNotFound
		}

		return repository.File{}, fmt.Errorf(
			"failed to get file: %w",
			err,
		)
	}

	if time.Now().After(file.ExpiresAt) {
		return repository.File{}, ErrFileExpired
	}

	if _, err := os.Stat(file.FilePath); err != nil {
		if os.IsNotExist(err) {
			return repository.File{}, ErrFileNotFound
		}

		return repository.File{}, fmt.Errorf(
			"failed to check file: %w",
			err,
		)
	}

	return file, nil
}

func (s *FileService) CleanupExpiredFiles(
	ctx context.Context,
) error {

	files, err := s.repo.GetExpiredFiles(ctx)
	if err != nil {
		return fmt.Errorf(
			"failed to get expired files: %w",
			err,
		)
	}

	for _, file := range files {

		err := os.Remove(file.FilePath)

		if err != nil && !os.IsNotExist(err) {
			fmt.Printf(
				"failed to delete file %s: %v\n",
				file.ID,
				err,
			)

			continue
		}

		if err := s.repo.Delete(ctx, file.ID); err != nil {
			fmt.Printf(
				"failed to delete database record %s: %v\n",
				file.ID,
				err,
			)

			continue
		}

		fmt.Printf(
			"deleted expired file: %s\n",
			file.ID,
		)
	}

	return nil
}