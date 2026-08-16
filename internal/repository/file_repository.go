package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type File struct {
	ID           string
	OriginalName string
	StoredName   string
	FilePath     string
	FileSize     int64
	ExpiresAt    string
}

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) Create(ctx context.Context, file File) error {
	query := `
		INSERT INTO files (
			id,
			original_name,
			stored_name,
			file_path,
			file_size,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		file.ID,
		file.OriginalName,
		file.StoredName,
		file.FilePath,
		file.FileSize,
		file.ExpiresAt,
	)

	return err
}