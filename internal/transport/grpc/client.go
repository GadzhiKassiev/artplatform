package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "artplatform/backend/proto/auth"
	coursepb "artplatform/backend/proto/course"
)

func NewAuthClient(addr string) (authpb.AuthServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return authpb.NewAuthServiceClient(conn), conn, nil
}

func NewCourseClient(addr string) (coursepb.CourseServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return coursepb.NewCourseServiceClient(conn), conn, nil
}
