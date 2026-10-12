package grpc

import (
	"google.golang.org/grpc"

	mediapb "artplatform/backend/proto/media"
)

type MediaClient struct {
	Client mediapb.MediaServiceClient
	conn   *grpc.ClientConn
}

func NewMediaClient(addr string) (*MediaClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}
	return &MediaClient{
			Client: mediapb.NewMediaServiceClient(conn),
			conn:   conn,
		},
		nil
}

func (c *MediaClient) Close() error {
	return c.conn.Close()
}
