package grpc

import (
	"google.golang.org/grpc"

	coursepb "artplatform/backend/proto/course"
)

type CourseClient struct {
	Client coursepb.CourseServiceClient
	conn   *grpc.ClientConn
}

func NewCourseClient(addr string) (*CourseClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}
	return &CourseClient{
			Client: coursepb.NewCourseServiceClient(conn),
			conn:   conn,
		},
		nil
}

func (c *CourseClient) Close() error {
	return c.conn.Close()
}
