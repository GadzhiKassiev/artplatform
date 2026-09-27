package model

import "time"

type CourseID = string

type Status string

const (
	StatusDraft     Status = "DRAFT"
	StatusPublished Status = "PUBLISHED"
)

type Course struct {
	ID          CourseID
	AuthorID    string
	Title       string
	Description string
	Price       float64
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
