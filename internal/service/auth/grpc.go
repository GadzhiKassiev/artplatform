package auth

import (
	"context"

	authpb "artplatform/backend/proto/auth"
)

type GRPCServer struct {
	authpb.UnimplementedAuthServiceServer
	service *Service
}

var _ authpb.AuthServiceServer = (*GRPCServer)(nil)

func NewGRPCServer(svc *Service) *GRPCServer {
	return &GRPCServer{service: svc}
}

func (s *GRPCServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	result, err := s.service.Register(ctx, RegisterInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if err != nil {
		return nil, err
	}
	return &authpb.RegisterResponse{
		UserId: result.UserID,
		Token:  result.Token,
	}, nil
}

func (s *GRPCServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	result, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return &authpb.LoginResponse{
		UserId: result.UserID,
		Token:  result.Token,
	}, nil
}

func (s *GRPCServer) ValidateToken(ctx context.Context, req *authpb.ValidateTokenRequest) (*authpb.ValidateTokenResponse, error) {
	userID, role, err := s.service.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &authpb.ValidateTokenResponse{
		UserId: userID,
		Role:   role,
	}, nil
}

func (s *GRPCServer) GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.GetUserResponse, error) {
	user, err := s.service.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &authpb.GetUserResponse{
		UserId: user.ID,
		Email:  user.Email,
		Role:   string(user.Role),
	}, nil
}
