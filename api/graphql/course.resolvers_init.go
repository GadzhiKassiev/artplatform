package graphql

import (
	"artplatform/backend/api/graphql/generated/model"
	coursepb "artplatform/backend/proto/course"
)

func toCourseModel(c *coursepb.Course) *model.Course {
	return &model.Course{
		ID:          c.Id,
		AuthorID:    c.AuthorId,
		Title:       c.Title,
		Description: c.Description,
		Price:       c.Price,
		Status:      c.Status,
	}
}
