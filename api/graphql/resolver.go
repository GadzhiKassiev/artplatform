package graphql

import (
	authpb "artplatform/backend/proto/auth"
	coursepb "artplatform/backend/proto/course"
)

type Resolver struct {
	AuthClient   authpb.AuthServiceClient
	CourseClient coursepb.CourseServiceClient
}
