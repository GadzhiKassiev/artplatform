package course

import (
	"context"

	"artplatform/backend/internal/service/course/model"
	grpctransport "artplatform/backend/internal/transport/grpc"
	coursepb "artplatform/backend/proto/course"
)

type GRPCServer struct {
	coursepb.UnimplementedCourseServiceServer
	service *Service
}

var _ coursepb.CourseServiceServer = (*GRPCServer)(nil)

func NewGRPCServer(svc *Service) *GRPCServer {
	return &GRPCServer{service: svc}
}

func (s *GRPCServer) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CourseResponse, error) {
	c, err := s.service.Create(ctx, CreateInput{
		AuthorID:    req.AuthorId,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &coursepb.CourseResponse{Course: toProto(c)}, nil
}

func (s *GRPCServer) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.CourseResponse, error) {
	c, err := s.service.GetByID(ctx, req.Id)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &coursepb.CourseResponse{Course: toProto(c)}, nil
}

func (s *GRPCServer) GetCourses(ctx context.Context, req *coursepb.GetCoursesRequest) (*coursepb.GetCoursesResponse, error) {
	var courses []model.Course
	var err error
	if req.AuthorId != "" {
		courses, err = s.service.GetByAuthor(ctx, req.AuthorId)
	} else {
		courses, err = s.service.GetAll(ctx)
	}
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	resp := &coursepb.GetCoursesResponse{Courses: make([]*coursepb.Course, len(courses))}
	for i, c := range courses {
		resp.Courses[i] = toProto(c)
	}
	return resp, nil
}

func (s *GRPCServer) UpdateCourse(ctx context.Context, req *coursepb.UpdateCourseRequest) (*coursepb.CourseResponse, error) {
	c, err := s.service.Update(ctx, UpdateInput{
		ID:          req.Id,
		RequesterID: req.RequesterId,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Status:      model.Status(req.Status),
	})
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &coursepb.CourseResponse{Course: toProto(c)}, nil
}

func (s *GRPCServer) DeleteCourse(ctx context.Context, req *coursepb.DeleteCourseRequest) (*coursepb.DeleteCourseResponse, error) {
	err := s.service.Delete(ctx, req.Id, req.RequesterId)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &coursepb.DeleteCourseResponse{Success: true}, nil
}

func toProto(c model.Course) *coursepb.Course {
	return &coursepb.Course{
		Id:          c.ID,
		AuthorId:    c.AuthorID,
		Title:       c.Title,
		Description: c.Description,
		Price:       c.Price,
		Status:      string(c.Status),
	}
}

func (s *GRPCServer) PublishCourse(ctx context.Context, req *coursepb.PublishCourseRequest) (*coursepb.CourseResponse, error) {
	c, err := s.service.Publish(ctx, req.Id, req.RequesterId)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &coursepb.CourseResponse{Course: toProto(c)}, nil
}
