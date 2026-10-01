package media

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"artplatform/backend/internal/logging"
	mediaerr "artplatform/backend/internal/service/media/err"
	"artplatform/backend/internal/service/media/model"
	"artplatform/backend/internal/service/media/storage"
	"artplatform/backend/internal/transport/kafka"
)

const (
	TopicMediaUploaded = "media.uploaded"
	PresignedTTL       = 1 * time.Hour
)

type Service struct {
	storage  storage.Storage
	minio    *storage.MinIOClient
	producer *kafka.Producer
}

func New(st storage.Storage, minio *storage.MinIOClient, producer *kafka.Producer) *Service {
	return &Service{
		storage:  st,
		minio:    minio,
		producer: producer,
	}
}

type UploadInput struct {
	OwnerID     string
	CourseID    string
	FileName    string
	ContentType string
	Content     []byte
}

type MediaUploadedEvent struct {
	MediaID     string `json:"media_id"`
	OriginalKey string `json:"original_key"`
	ContentType string `json:"content_type"`
}

func (s *Service) Upload(ctx context.Context, input UploadInput) (model.MediaFile, error) {
	logging.ContextInfo(ctx, "uploading media",
		logging.NewKV("ownerID", input.OwnerID),
		logging.NewKV("fileName", input.FileName),
		logging.NewKV("size", len(input.Content)),
	)

	mediaID := uuid.NewString()
	originalKey := fmt.Sprintf("media/%s/original/%s", mediaID, input.FileName)

	if err := s.minio.Upload(ctx, originalKey, input.Content, input.ContentType); err != nil {
		logging.ContextErrorE(ctx, "failed to upload to minio", err)
		return model.MediaFile{}, mediaerr.ErrStorageError
	}

	media := model.MediaFile{
		ID:          mediaID,
		OwnerID:     input.OwnerID,
		CourseID:    input.CourseID,
		FileName:    input.FileName,
		ContentType: input.ContentType,
		Size:        int64(len(input.Content)),
		OriginalKey: originalKey,
		Status:      model.StatusProcessing,
	}

	created, err := s.storage.CreateMedia(ctx, media)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to save media metadata", err)
		return model.MediaFile{}, err
	}

	event := MediaUploadedEvent{
		MediaID:     created.ID,
		OriginalKey: created.OriginalKey,
		ContentType: created.ContentType,
	}
	if err := s.producer.Publish(ctx, TopicMediaUploaded, created.ID, event); err != nil {
		logging.ContextErrorE(ctx, "failed to publish media.uploaded", err)
	} else {
		logging.ContextInfo(ctx, "media.uploaded published", logging.NewKV("mediaID", created.ID))
	}

	logging.ContextInfo(ctx, "media uploaded", logging.NewKV("mediaID", created.ID))
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id model.MediaID) (model.MediaFile, error) {
	return s.storage.GetMediaByID(ctx, id)
}

func (s *Service) GetPresignedURL(ctx context.Context, id model.MediaID, requesterID string) (string, string, error) {
	media, err := s.storage.GetMediaByID(ctx, id)
	if err != nil {
		return "", "", err
	}

	if media.OwnerID != requesterID {
		// For prototype: only owner can get URL
		// In real system: also allow course students, admins, etc.
		logging.ContextWarn(ctx, "access denied for presigned URL",
			logging.NewKV("mediaID", id),
			logging.NewKV("requesterID", requesterID),
		)
		return "", "", mediaerr.ErrAccessDenied
	}

	key := media.OriginalKey
	kind := "original"
	if media.PreviewKey != "" && media.Status == model.StatusReady {
		key = media.PreviewKey
		kind = "preview"
	}

	url, err := s.minio.PresignedGetURL(ctx, key, PresignedTTL)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to generate presigned URL", err)
		return "", "", mediaerr.ErrStorageError
	}

	logging.ContextInfo(ctx, "presigned URL generated",
		logging.NewKV("mediaID", id),
		logging.NewKV("kind", kind),
	)
	return url, kind, nil
}
