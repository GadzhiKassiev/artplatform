package grpc

import (
	"google.golang.org/grpc"

	authpb "artplatform/backend/proto/auth"
)

type AuthClient struct {
	Client authpb.AuthServiceClient
	conn   *grpc.ClientConn
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}
	return &AuthClient{Client: authpb.NewAuthServiceClient(conn),
			conn: conn,
		},
		nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}
