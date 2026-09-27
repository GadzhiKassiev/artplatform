package storage

import (
	"context"

	"artplatform/backend/internal/service/course/model"
)

type Storage interface {
	CreateCourse(ctx context.Context, course model.Course) (model.Course, error)
	GetCourseByID(ctx context.Context, id model.CourseID) (model.Course, error)
	GetCoursesByAuthor(ctx context.Context, authorID string) ([]model.Course, error)
	GetAllCourses(ctx context.Context) ([]model.Course, error)
	UpdateCourse(ctx context.Context, course model.Course) (model.Course, error)
	DeleteCourse(ctx context.Context, id model.CourseID) error
}
