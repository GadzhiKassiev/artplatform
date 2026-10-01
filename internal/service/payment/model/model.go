package model

import "time"

type TransactionID = string

type Status string

const (
	StatusPending Status = "PENDING"
	StatusSuccess Status = "SUCCESS"
	StatusFailed  Status = "FAILED"
)

type Transaction struct {
	ID        TransactionID
	UserID    string
	CourseID  string
	Amount    float64
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
