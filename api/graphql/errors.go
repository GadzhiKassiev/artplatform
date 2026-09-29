package graphql

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"

	"artplatform/backend/internal/errs"
	grpctransport "artplatform/backend/internal/transport/grpc"
)

func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	err = grpctransport.DecodeError(err)

	var domainErr *errs.Error
	if errors.As(err, &domainErr) {
		return &gqlerror.Error{
			Message: domainErr.Message,
			Extensions: map[string]any{
				"code": string(domainErr.Code),
			},
			Path: graphql.GetPath(ctx),
		}
	}

	return &gqlerror.Error{
		Message: "internal error",
		Extensions: map[string]any{
			"code": string(errs.CodeInternal),
		},
		Path: graphql.GetPath(ctx),
	}
}
