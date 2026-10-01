package model

import "time"

type CourseStats struct {
	CourseID      string
	ViewCount     int64
	PurchaseCount int64
	UpdatedAt     time.Time
}
