package storage

import (
	"context"

	"artplatform/backend/internal/service/media/model"
)

type Storage interface {
	CreateMedia(ctx context.Context, media model.MediaFile) (model.MediaFile, error)
	GetMediaByID(ctx context.Context, id model.MediaID) (model.MediaFile, error)
	UpdateMedia(ctx context.Context, media model.MediaFile) (model.MediaFile, error)
	UpdateStatus(ctx context.Context, id model.MediaID, status model.Status, previewKey *string) error
}
