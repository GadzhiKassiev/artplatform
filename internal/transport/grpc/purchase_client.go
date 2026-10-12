package grpc

import (
	"google.golang.org/grpc"

	purchasepb "artplatform/backend/proto/purchase"
)

type PurchaseClient struct {
	Client purchasepb.PurchaseServiceClient
	conn   *grpc.ClientConn
}

func NewPurchaseClient(addr string) (*PurchaseClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}
	return &PurchaseClient{
			Client: purchasepb.NewPurchaseServiceClient(conn),
			conn:   conn,
		},
		nil
}

func (c *PurchaseClient) Close() error {
	return c.conn.Close()
}
