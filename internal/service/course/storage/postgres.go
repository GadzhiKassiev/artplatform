package storage

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"artplatform/backend/internal/service/course/model"
)

var ErrCourseNotFound = errors.New("course not found")

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgres(pool *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) CreateCourse(ctx context.Context, course model.Course) (model.Course, error) {
	now := time.Now().UTC()
	course.CreatedAt = now
	course.UpdatedAt = now

	query := `
		INSERT INTO courses (id, author_id, title, description, price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := s.pool.Exec(ctx, query,
		course.ID, course.AuthorID, course.Title, course.Description,
		course.Price, course.Status, course.CreatedAt, course.UpdatedAt,
	)
	if err != nil {
		return model.Course{}, err
	}
	return course, nil
}

func (s *PostgresStorage) GetCourseByID(ctx context.Context, id model.CourseID) (model.Course, error) {
	query := `
		SELECT id, author_id, title, description, price, status, created_at, updated_at
		FROM courses
		WHERE id = $1
	`
	var course model.Course
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&course.ID, &course.AuthorID, &course.Title, &course.Description,
		&course.Price, &course.Status, &course.CreatedAt, &course.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Course{}, ErrCourseNotFound
	}
	if err != nil {
		return model.Course{}, err
	}
	return course, nil
}

func (s *PostgresStorage) GetCoursesByAuthor(ctx context.Context, authorID string) ([]model.Course, error) {
	query := `
		SELECT id, author_id, title, description, price, status, created_at, updated_at
		FROM courses
		WHERE author_id = $1
		ORDER BY created_at DESC
	`
	return s.scanCourses(ctx, query, authorID)
}

func (s *PostgresStorage) GetAllCourses(ctx context.Context) ([]model.Course, error) {
	query := `
		SELECT id, author_id, title, description, price, status, created_at, updated_at
		FROM courses
		WHERE status = 'PUBLISHED'
		ORDER BY created_at DESC
	`
	return s.scanCourses(ctx, query)
}

func (s *PostgresStorage) UpdateCourse(ctx context.Context, course model.Course) (model.Course, error) {
	course.UpdatedAt = time.Now().UTC()

	query := `
		UPDATE courses
		SET title = $2, description = $3, price = $4, status = $5, updated_at = $6
		WHERE id = $1
	`
	tag, err := s.pool.Exec(ctx, query,
		course.ID, course.Title, course.Description, course.Price, course.Status, course.UpdatedAt,
	)
	if err != nil {
		return model.Course{}, err
	}
	if tag.RowsAffected() == 0 {
		return model.Course{}, ErrCourseNotFound
	}
	return course, nil
}

func (s *PostgresStorage) DeleteCourse(ctx context.Context, id model.CourseID) error {
	query := `DELETE FROM courses WHERE id = $1`
	tag, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrCourseNotFound
	}
	return nil
}

func (s *PostgresStorage) scanCourses(ctx context.Context, query string, args ...any) ([]model.Course, error) {
	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []model.Course
	for rows.Next() {
		var c model.Course
		if err := rows.Scan(
			&c.ID, &c.AuthorID, &c.Title, &c.Description,
			&c.Price, &c.Status, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		courses = append(courses, c)
	}
	return courses, rows.Err()
}
