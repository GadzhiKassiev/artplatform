package grpc

import (
	"google.golang.org/grpc"

	activitypb "artplatform/backend/proto/activity"
)

type ActivityClient struct {
	Client activitypb.ActivityServiceClient
	conn   *grpc.ClientConn
}

func NewActivityClient(addr string) (*ActivityClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}
	return &ActivityClient{
			Client: activitypb.NewActivityServiceClient(conn),
			conn:   conn,
		},
		nil
}

func (c *ActivityClient) Close() error {
	return c.conn.Close()
}
