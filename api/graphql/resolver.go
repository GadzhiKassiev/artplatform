package graphql

import (
	activitypb "artplatform/backend/proto/activity"
	authpb "artplatform/backend/proto/auth"
	coursepb "artplatform/backend/proto/course"
	mediapb "artplatform/backend/proto/media"
	paymentpb "artplatform/backend/proto/payment"
	purchasepb "artplatform/backend/proto/purchase"
)

type Resolver struct {
	ActivityClient activitypb.ActivityServiceClient
	AuthClient     authpb.AuthServiceClient
	CourseClient   coursepb.CourseServiceClient
	MediaClient    mediapb.MediaServiceClient
	PaymentClient  paymentpb.PaymentServiceClient
	PurchaseClient purchasepb.PurchaseServiceClient
}
