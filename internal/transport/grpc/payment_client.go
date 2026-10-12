package grpc

import (
	"google.golang.org/grpc"

	paymentpb "artplatform/backend/proto/payment"
)

type PaymentClient struct {
	Client paymentpb.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentClient(addr string) (*PaymentClient, error) {
	conn, err := newConn(addr)
	if err != nil {
		return nil, err
	}
	return &PaymentClient{
			Client: paymentpb.NewPaymentServiceClient(conn),
			conn:   conn,
		},
		nil
}

func (c *PaymentClient) Close() error {
	return c.conn.Close()
}
