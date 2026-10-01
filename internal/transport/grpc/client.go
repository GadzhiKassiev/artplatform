package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	activitypb "artplatform/backend/proto/activity"
	authpb "artplatform/backend/proto/auth"
	coursepb "artplatform/backend/proto/course"
	mediapb "artplatform/backend/proto/media"
	paymentpb "artplatform/backend/proto/payment"
	purchasepb "artplatform/backend/proto/purchase"
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

func NewMediaClient(addr string) (mediapb.MediaServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return mediapb.NewMediaServiceClient(conn), conn, nil
}

func NewPaymentClient(addr string) (paymentpb.PaymentServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return paymentpb.NewPaymentServiceClient(conn), conn, nil
}

func NewPurchaseClient(addr string) (purchasepb.PurchaseServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return purchasepb.NewPurchaseServiceClient(conn), conn, nil
}

func NewActivityClient(addr string) (activitypb.ActivityServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	return activitypb.NewActivityServiceClient(conn), conn, nil
}
