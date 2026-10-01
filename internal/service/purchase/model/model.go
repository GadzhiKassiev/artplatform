package model

import "time"

type PurchaseID = string

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	StatusRefunded  Status = "REFUNDED"
	StatusCancelled Status = "CANCELLED"
)

type Purchase struct {
	ID            PurchaseID
	UserID        string
	CourseID      string
	Amount        float64
	Status        Status
	TransactionID string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
