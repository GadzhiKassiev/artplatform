package payment

import (
	"context"

	"artplatform/backend/internal/service/payment/model"
	grpctransport "artplatform/backend/internal/transport/grpc"
	paymentpb "artplatform/backend/proto/payment"
)

type GRPCServer struct {
	paymentpb.UnimplementedPaymentServiceServer
	service *Service
}

var _ paymentpb.PaymentServiceServer = (*GRPCServer)(nil)

func NewGRPCServer(svc *Service) *GRPCServer {
	return &GRPCServer{service: svc}
}

func (s *GRPCServer) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.TransactionResponse, error) {
	tx, err := s.service.ProcessPayment(ctx, ProcessInput{
		UserID:   req.UserId,
		CourseID: req.CourseId,
		Amount:   req.Amount,
	})
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &paymentpb.TransactionResponse{Transaction: toProto(tx)}, nil
}

func (s *GRPCServer) GetTransaction(ctx context.Context, req *paymentpb.GetTransactionRequest) (*paymentpb.TransactionResponse, error) {
	tx, err := s.service.GetTransaction(ctx, req.Id)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &paymentpb.TransactionResponse{Transaction: toProto(tx)}, nil
}

func toProto(tx model.Transaction) *paymentpb.Transaction {
	return &paymentpb.Transaction{
		Id:       tx.ID,
		UserId:   tx.UserID,
		CourseId: tx.CourseID,
		Amount:   tx.Amount,
		Status:   string(tx.Status),
	}
}
