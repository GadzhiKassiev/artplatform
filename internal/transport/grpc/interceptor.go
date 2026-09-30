package grpc

import (
	"context"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"

	"artplatform/backend/internal/errs"
)

// ValidationInterceptor validates incoming gRPC requests using Protovalidate.
func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if msg, ok := req.(protoreflect.ProtoMessage); ok {
			if err := protovalidate.Validate(msg); err != nil {
				return nil, EncodeError(errs.New(errs.CodeInvalidInput, err.Error()))
			}
		}
		return handler(ctx, req)
	}
}
