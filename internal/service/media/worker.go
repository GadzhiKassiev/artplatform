package media

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/disintegration/imaging"

	"artplatform/backend/internal/logging"
	"artplatform/backend/internal/service/media/model"
	"artplatform/backend/internal/service/media/storage"
	"artplatform/backend/internal/transport/kafka"
)

const (
	PreviewWidth   = 400
	PreviewQuality = 85
	MaxFileSize    = 50 * 1024 * 1024
)

type Worker struct {
	storage  storage.Storage
	minio    *storage.MinIOClient
	handler  kafka.MessageHandler
	consumer *kafka.Consumer
	timeout  time.Duration
}

func NewWorker(st storage.Storage, minio *storage.MinIOClient, brokers []string, timeout time.Duration) *Worker {
	w := &Worker{
		storage: st,
		minio:   minio,
		timeout: timeout,
	}

	w.handler = w.processMessage
	w.consumer = kafka.NewConsumer(brokers, TopicMediaUploaded, "media-worker", w.handler)

	return w
}

func (w *Worker) Run(ctx context.Context) error {
	logging.ContextInfo(ctx, "media-worker started, listening for media.uploaded")
	return w.consumer.Run(ctx)
}

func (w *Worker) Close() error {
	return w.consumer.Close()
}

// processMessage handles a single media.uploaded event.
func (w *Worker) processMessage(ctx context.Context, key string, value []byte) error {
	event, err := kafka.Decode[MediaUploadedEvent](value)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to decode media.uploaded event", err)
		return err
	}

	acquired, err := w.storage.TryAcquireForProcessing(ctx, event.MediaID, w.timeout)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to acquire media for processing", err)
		return err
	}
	if !acquired {
		media, err := w.storage.GetMediaByID(ctx, event.MediaID)
		if err != nil {
			return err
		}
		if media.Status == model.StatusReady {
			return nil
		}
		return errors.New("media is being processed by another worker")
	}

	if event.Size > MaxFileSize {
		logging.ContextWarn(ctx, "file too large, skipping",
			logging.NewKV("mediaID", event.MediaID),
			logging.NewKV("size", event.Size),
		)
		_ = w.storage.UpdateStatus(ctx, event.MediaID, model.StatusFailed, nil)
		return nil
	}

	logging.ContextInfo(ctx, "processing media",
		logging.NewKV("mediaID", event.MediaID),
		logging.NewKV("contentType", event.ContentType),
	)

	original, err := w.minio.Download(ctx, event.OriginalKey)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to download original", err)
		_ = w.storage.UpdateStatus(ctx, event.MediaID, model.StatusFailed, nil)
		return err
	}

	preview, previewKey, err := w.generatePreview(event, original)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to generate preview", err)
		_ = w.storage.UpdateStatus(ctx, event.MediaID, model.StatusFailed, nil)
		return err
	}

	if err := w.minio.Upload(ctx, previewKey, preview, "image/jpeg"); err != nil {
		logging.ContextErrorE(ctx, "failed to upload preview", err)
		_ = w.storage.UpdateStatus(ctx, event.MediaID, model.StatusFailed, nil)
		return err
	}

	if err := w.storage.UpdateStatus(ctx, event.MediaID, model.StatusReady, &previewKey); err != nil {
		logging.ContextErrorE(ctx, "failed to update status", err)
		return err
	}

	logging.ContextInfo(ctx, "preview generated",
		logging.NewKV("mediaID", event.MediaID),
		logging.NewKV("previewKey", previewKey),
	)
	return nil
}

func (w *Worker) generatePreview(event MediaUploadedEvent, original []byte) ([]byte, string, error) {
	img, err := imaging.Decode(bytes.NewReader(original))
	if err != nil {
		return nil, "", fmt.Errorf("decode image: %w", err)
	}

	resized := imaging.Resize(img, PreviewWidth, 0, imaging.Lanczos)

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, resized, imaging.JPEG, imaging.JPEGQuality(PreviewQuality)); err != nil {
		return nil, "", fmt.Errorf("encode preview: %w", err)
	}

	previewKey := strings.Replace(event.OriginalKey, "/original/", "/preview/", 1)
	if idx := strings.LastIndex(previewKey, "."); idx != -1 {
		previewKey = previewKey[:idx] + ".jpg"
	}

	return buf.Bytes(), previewKey, nil
}
