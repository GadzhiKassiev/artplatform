package model

import "time"

type MediaID = string

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusProcessing Status = "PROCESSING"
	StatusReady      Status = "READY"
	StatusFailed     Status = "FAILED"
)

type MediaFile struct {
	ID          MediaID
	OwnerID     string
	CourseID    *string
	FileName    string
	ContentType string
	Size        int64
	OriginalKey string
	PreviewKey  *string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
