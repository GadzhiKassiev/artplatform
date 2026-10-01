package graphql

import (
	"artplatform/backend/api/graphql/generated/model"
	mediapb "artplatform/backend/proto/media"
)

func toMediaModel(m *mediapb.MediaFile) *model.Media {
	if m == nil {
		return nil
	}
	return &model.Media{
		ID:          m.Id,
		OwnerID:     m.OwnerId,
		CourseID:    strPtr(m.CourseId),
		FileName:    m.FileName,
		ContentType: m.ContentType,
		Size:        int(m.Size),
		Status:      m.Status,
		PreviewKey:  strPtr(m.PreviewKey),
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
