package graphql

import (
	"context"

	"artplatform/backend/api/graphql/generated/model"
	activitypb "artplatform/backend/proto/activity"
)

func (r *queryResolver) GetCourseStats(ctx context.Context, courseID string) (*model.CourseStats, error) {
	resp, err := r.ActivityClient.GetCourseStats(ctx, &activitypb.GetCourseStatsRequest{CourseId: courseID})
	if err != nil {
		return nil, err
	}
	return &model.CourseStats{
		CourseID:      resp.Stats.CourseId,
		ViewCount:     int(resp.Stats.ViewCount),
		PurchaseCount: int(resp.Stats.PurchaseCount),
	}, nil
}

func (r *mutationResolver) RecordView(ctx context.Context, courseID string) (bool, error) {
	resp, err := r.ActivityClient.RecordView(ctx, &activitypb.RecordViewRequest{CourseId: courseID})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}
