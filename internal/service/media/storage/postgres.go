package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"artplatform/backend/internal/service/media/model"
)

var ErrMediaNotFound = errors.New("media not found")

type PostgresStorage struct {
	pool *pgxpool.Pool
}

var _ Storage = (*PostgresStorage)(nil)

func NewPostgres(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) CreateMedia(ctx context.Context, media model.MediaFile) (model.MediaFile, error) {
	now := time.Now().UTC()
	media.CreatedAt = now
	media.UpdatedAt = now

	query := `
		INSERT INTO media_files (id, owner_id, course_id, file_name, content_type, size,
		                         original_key, preview_key, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := s.pool.Exec(ctx, query,
		media.ID, media.OwnerID, media.CourseID, media.FileName, media.ContentType,
		media.Size, media.OriginalKey, media.PreviewKey, media.Status,
		media.CreatedAt, media.UpdatedAt,
	)
	if err != nil {
		return model.MediaFile{}, err
	}
	return media, nil
}

func (s *PostgresStorage) GetMediaByID(ctx context.Context, id model.MediaID) (model.MediaFile, error) {
	query := `
		SELECT id, owner_id, course_id, file_name, content_type, size,
		       original_key, preview_key, status, created_at, updated_at
		FROM media_files
		WHERE id = $1
	`
	var m model.MediaFile
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.OwnerID, &m.CourseID, &m.FileName, &m.ContentType, &m.Size,
		&m.OriginalKey, &m.PreviewKey, &m.Status, &m.CreatedAt, &m.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.MediaFile{}, ErrMediaNotFound
	}
	if err != nil {
		return model.MediaFile{}, err
	}
	return m, nil
}

func (s *PostgresStorage) UpdateMedia(ctx context.Context, media model.MediaFile) (model.MediaFile, error) {
	media.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE media_files
		SET preview_key = $2, status = $3, updated_at = $4
		WHERE id = $1
	`
	tag, err := s.pool.Exec(ctx, query,
		media.ID, media.PreviewKey, media.Status, media.UpdatedAt,
	)
	if err != nil {
		return model.MediaFile{}, err
	}
	if tag.RowsAffected() == 0 {
		return model.MediaFile{}, ErrMediaNotFound
	}
	return media, nil
}

func (s *PostgresStorage) UpdateStatus(
	ctx context.Context,
	id model.MediaID,
	status model.Status,
	previewKey *string,
) error {
	query := `
		UPDATE media_files
		SET status = $2, preview_key = $3, updated_at = now()
		WHERE id = $1
	`
	tag, err := s.pool.Exec(ctx, query, id, status, previewKey)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrMediaNotFound
	}
	return nil
}
