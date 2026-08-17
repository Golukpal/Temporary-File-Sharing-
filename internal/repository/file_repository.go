package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrFileNotFound = errors.New("file not found")

type File struct {
	ID           string
	OriginalName string
	StoredName   string
	FilePath     string
	FileSize     int64
	ExpiresAt    time.Time
}

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(db *pgxpool.Pool) *FileRepository {
	return &FileRepository{
		db: db,
	}
}

func (r *FileRepository) Create(
	ctx context.Context,
	file File,
) error {

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

func (r *FileRepository) GetByID(
	ctx context.Context,
	id string,
) (File, error) {

	query := `
		SELECT
			id,
			original_name,
			stored_name,
			file_path,
			file_size,
			expires_at
		FROM files
		WHERE id = $1
	`

	var file File

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&file.ID,
		&file.OriginalName,
		&file.StoredName,
		&file.FilePath,
		&file.FileSize,
		&file.ExpiresAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return File{}, ErrFileNotFound
	}

	if err != nil {
		return File{}, err
	}

	return file, nil
}

func (r *FileRepository) GetExpiredFiles(
	ctx context.Context,
) ([]File, error) {

	query := `
		SELECT
			id,
			original_name,
			stored_name,
			file_path,
			file_size,
			expires_at
		FROM files
		WHERE expires_at <= NOW()
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var files []File

	for rows.Next() {
		var file File

		err := rows.Scan(
			&file.ID,
			&file.OriginalName,
			&file.StoredName,
			&file.FilePath,
			&file.FileSize,
			&file.ExpiresAt,
		)

		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *FileRepository) Delete(
	ctx context.Context,
	id string,
) error {

	query := `
		DELETE FROM files
		WHERE id = $1
	`

	_, err := r.db.Exec(ctx, query, id)

	return err
}