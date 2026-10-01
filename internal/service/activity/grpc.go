package activity

import (
	"context"

	"artplatform/backend/internal/service/activity/model"
	grpctransport "artplatform/backend/internal/transport/grpc"
	activitypb "artplatform/backend/proto/activity"
)

type GRPCServer struct {
	activitypb.UnimplementedActivityServiceServer
	service *Service
}

var _ activitypb.ActivityServiceServer = (*GRPCServer)(nil)

func NewGRPCServer(svc *Service) *GRPCServer {
	return &GRPCServer{service: svc}
}

func (s *GRPCServer) GetCourseStats(ctx context.Context, req *activitypb.GetCourseStatsRequest) (*activitypb.CourseStatsResponse, error) {
	stats, err := s.service.GetStats(ctx, req.CourseId)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &activitypb.CourseStatsResponse{Stats: toProto(stats)}, nil
}

func (s *GRPCServer) RecordView(ctx context.Context, req *activitypb.RecordViewRequest) (*activitypb.RecordViewResponse, error) {
	if err := s.service.RecordView(ctx, req.CourseId); err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &activitypb.RecordViewResponse{Success: true}, nil
}

func toProto(s model.CourseStats) *activitypb.CourseStats {
	return &activitypb.CourseStats{
		CourseId:      s.CourseID,
		ViewCount:     s.ViewCount,
		PurchaseCount: s.PurchaseCount,
	}
}
