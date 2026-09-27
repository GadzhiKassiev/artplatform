package graphql

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "userID"

var ErrUnauthenticated = errors.New("unauthenticated")

func withUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func userIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok || userID == "" {
		return "", ErrUnauthenticated
	}
	return userID, nil
}
