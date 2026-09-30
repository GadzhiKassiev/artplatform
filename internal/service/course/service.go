package course

import (
	"context"

	"github.com/google/uuid"

	"artplatform/backend/internal/logging"
	courseerr "artplatform/backend/internal/service/course/err"
	"artplatform/backend/internal/service/course/model"
	"artplatform/backend/internal/service/course/storage"
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
	logging.ContextInfo(ctx, "creating course",
		logging.NewKV("authorID", input.AuthorID),
		logging.NewKV("title", input.Title),
	)

	course := model.Course{
		ID:          uuid.NewString(),
		AuthorID:    input.AuthorID,
		Title:       input.Title,
		Description: input.Description,
		Price:       input.Price,
		Status:      model.StatusDraft,
	}

	created, err := s.storage.CreateCourse(ctx, course)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to create course", err)
		return model.Course{}, err
	}

	logging.ContextInfo(ctx, "course created", logging.NewKV("courseID", created.ID))
	return created, nil
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
	logging.ContextInfo(ctx, "updating course",
		logging.NewKV("courseID", input.ID),
		logging.NewKV("requesterID", input.RequesterID),
	)

	course, err := s.storage.GetCourseByID(ctx, input.ID)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to get course", err)
		return model.Course{}, err
	}

	if course.AuthorID != input.RequesterID {
		logging.ContextWarn(ctx, "not owner",
			logging.NewKV("courseID", input.ID),
			logging.NewKV("requesterID", input.RequesterID),
		)
		return model.Course{}, courseerr.ErrNotOwner
	}

	course.Title = input.Title
	course.Description = input.Description
	course.Price = input.Price
	if input.Status != "" {
		course.Status = input.Status
	}

	updated, err := s.storage.UpdateCourse(ctx, course)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to update course", err)
		return model.Course{}, err
	}

	logging.ContextInfo(ctx, "course updated", logging.NewKV("courseID", updated.ID))
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id model.CourseID, requesterID string) error {
	logging.ContextInfo(ctx, "deleting course",
		logging.NewKV("courseID", id),
		logging.NewKV("requesterID", requesterID),
	)

	course, err := s.storage.GetCourseByID(ctx, id)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to get course", err)
		return err
	}

	if course.AuthorID != requesterID {
		logging.ContextWarn(ctx, "not owner", logging.NewKV("courseID", id))
		return courseerr.ErrNotOwner
	}

	if err := s.storage.DeleteCourse(ctx, id); err != nil {
		logging.ContextErrorE(ctx, "failed to delete course", err)
		return err
	}

	logging.ContextInfo(ctx, "course deleted", logging.NewKV("courseID", id))
	return nil
}

func (s *Service) Publish(ctx context.Context, id model.CourseID, requesterID string) (model.Course, error) {
	logging.ContextInfo(ctx, "publishing course",
		logging.NewKV("courseID", id),
		logging.NewKV("requesterID", requesterID),
	)

	course, err := s.storage.GetCourseByID(ctx, id)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to get course", err)
		return model.Course{}, err
	}

	if course.AuthorID != requesterID {
		logging.ContextWarn(ctx, "not owner", logging.NewKV("courseID", id))
		return model.Course{}, courseerr.ErrNotOwner
	}

	if course.Status == model.StatusPublished {
		logging.ContextWarn(ctx, "course already published", logging.NewKV("courseID", id))
		return model.Course{}, courseerr.ErrAlreadyPublished
	}

	course.Status = model.StatusPublished
	updated, err := s.storage.UpdateCourse(ctx, course)
	if err != nil {
		logging.ContextErrorE(ctx, "failed to publish course", err)
		return model.Course{}, err
	}

	logging.ContextInfo(ctx, "course published", logging.NewKV("courseID", updated.ID))
	return updated, nil
}
