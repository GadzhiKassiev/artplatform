package course

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"artplatform/backend/internal/service/course/model"
	"artplatform/backend/internal/service/course/storage"
)

var (
	ErrInvalidPrice     = errors.New("price must be >= 0")
	ErrEmptyTitle       = errors.New("title must not be empty")
	ErrNotOwner         = errors.New("only author can modify course")
	ErrAlreadyPublished = errors.New("course already published")
)

type Service struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

type CreateInput struct {
	AuthorID    string
	Title       string
	Description string
	Price       float64
}

func (s *Service) Create(ctx context.Context, input CreateInput) (model.Course, error) {
	if input.Title == "" {
		return model.Course{}, ErrEmptyTitle
	}
	if input.Price < 0 {
		return model.Course{}, ErrInvalidPrice
	}

	course := model.Course{
		ID:          uuid.NewString(),
		AuthorID:    input.AuthorID,
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Status:      model.StatusDraft,
	}

	return s.storage.CreateCourse(ctx, course)
}

func (s *Service) GetByID(ctx context.Context, id model.CourseID) (model.Course, error) {
	return s.storage.GetCourseByID(ctx, id)
}

func (s *Service) GetByAuthor(ctx context.Context, authorID string) ([]model.Course, error) {
	return s.storage.GetCoursesByAuthor(ctx, authorID)
}

func (s *Service) GetAll(ctx context.Context) ([]model.Course, error) {
	return s.storage.GetAllCourses(ctx)
}

type UpdateInput struct {
	ID          model.CourseID
	RequesterID string
	Title       string
	Description string
	Price       float64
	Status      model.Status
}

func (s *Service) Update(ctx context.Context, input UpdateInput) (model.Course, error) {
	if input.Title == "" {
		return model.Course{}, ErrEmptyTitle
	}
	if input.Price < 0 {
		return model.Course{}, ErrInvalidPrice
	}

	course, err := s.storage.GetCourseByID(ctx, input.ID)
	if err != nil {
		return model.Course{}, err
	}

	if course.AuthorID != input.RequesterID {
		return model.Course{}, ErrNotOwner
	}

	course.Title = input.Title
	course.Description = input.Description
	course.Price = input.Price
	if input.Status != "" {
		course.Status = input.Status
	}

	return s.storage.UpdateCourse(ctx, course)
}

func (s *Service) Delete(ctx context.Context, id model.CourseID, requesterID string) error {
	course, err := s.storage.GetCourseByID(ctx, id)
	if err != nil {
		return err
	}

	if course.AuthorID != requesterID {
		return ErrNotOwner
	}

	return s.storage.DeleteCourse(ctx, id)
}

func (s *Service) Publish(ctx context.Context, id model.CourseID, requesterID string) (model.Course, error) {
	course, err := s.storage.GetCourseByID(ctx, id)
	if err != nil {
		return model.Course{}, err
	}

	if course.AuthorID != requesterID {
		return model.Course{}, ErrNotOwner
	}

	if course.Status == model.StatusPublished {
		return model.Course{}, ErrAlreadyPublished
	}

	course.Status = model.StatusPublished
	return s.storage.UpdateCourse(ctx, course)
}
