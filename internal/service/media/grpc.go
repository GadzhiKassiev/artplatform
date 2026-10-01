package media

import (
	"context"

	"artplatform/backend/internal/service/media/model"
	grpctransport "artplatform/backend/internal/transport/grpc"
	mediapb "artplatform/backend/proto/media"
)

type GRPCServer struct {
	mediapb.UnimplementedMediaServiceServer
	service *Service
}

var _ mediapb.MediaServiceServer = (*GRPCServer)(nil)

func NewGRPCServer(svc *Service) *GRPCServer {
	return &GRPCServer{service: svc}
}

func (s *GRPCServer) UploadMedia(ctx context.Context, req *mediapb.UploadMediaRequest) (*mediapb.MediaResponse, error) {
	media, err := s.service.Upload(ctx, UploadInput{
		OwnerID:     req.OwnerId,
		CourseID:    req.CourseId,
		FileName:    req.FileName,
		ContentType: req.ContentType,
		Content:     req.Content,
	})
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &mediapb.MediaResponse{Media: toProto(media)}, nil
}

func (s *GRPCServer) GetMedia(ctx context.Context, req *mediapb.GetMediaRequest) (*mediapb.MediaResponse, error) {
	media, err := s.service.GetByID(ctx, req.Id)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &mediapb.MediaResponse{Media: toProto(media)}, nil
}

func (s *GRPCServer) GetPresignedURL(ctx context.Context, req *mediapb.GetPresignedURLRequest) (*mediapb.GetPresignedURLResponse, error) {
	url, kind, err := s.service.GetPresignedURL(ctx, req.Id, req.RequesterId)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &mediapb.GetPresignedURLResponse{Url: url, Type: kind}, nil
}

func toProto(m model.MediaFile) *mediapb.MediaFile {
	return &mediapb.MediaFile{
		Id:          m.ID,
		OwnerId:     m.OwnerID,
		CourseId:    m.CourseID,
		FileName:    m.FileName,
		ContentType: m.ContentType,
		Size:        m.Size,
		OriginalKey: m.OriginalKey,
		PreviewKey:  m.PreviewKey,
		Status:      string(m.Status),
	}
}
