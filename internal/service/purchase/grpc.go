package purchase

import (
	"context"

	"artplatform/backend/internal/service/purchase/model"
	grpctransport "artplatform/backend/internal/transport/grpc"
	purchasepb "artplatform/backend/proto/purchase"
)

type GRPCServer struct {
	purchasepb.UnimplementedPurchaseServiceServer
	service *Service
}

var _ purchasepb.PurchaseServiceServer = (*GRPCServer)(nil)

func NewGRPCServer(svc *Service) *GRPCServer {
	return &GRPCServer{service: svc}
}

func (s *GRPCServer) CreatePurchase(ctx context.Context, req *purchasepb.CreatePurchaseRequest) (*purchasepb.PurchaseResponse, error) {
	p, err := s.service.Create(ctx, CreateInput{
		UserID:   req.UserId,
		CourseID: req.CourseId,
		Amount:   req.Amount,
	})
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &purchasepb.PurchaseResponse{Purchase: toProto(p)}, nil
}

func (s *GRPCServer) GetPurchase(ctx context.Context, req *purchasepb.GetPurchaseRequest) (*purchasepb.PurchaseResponse, error) {
	p, err := s.service.GetByID(ctx, req.Id)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &purchasepb.PurchaseResponse{Purchase: toProto(p)}, nil
}

func (s *GRPCServer) GetUserPurchases(ctx context.Context, req *purchasepb.GetUserPurchasesRequest) (*purchasepb.GetUserPurchasesResponse, error) {
	purchases, err := s.service.GetUserPurchases(ctx, req.UserId)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	resp := &purchasepb.GetUserPurchasesResponse{
		Purchases: make([]*purchasepb.Purchase, len(purchases)),
	}
	for i, p := range purchases {
		resp.Purchases[i] = toProto(p)
	}
	return resp, nil
}

func (s *GRPCServer) CheckAccess(ctx context.Context, req *purchasepb.CheckAccessRequest) (*purchasepb.CheckAccessResponse, error) {
	hasAccess, err := s.service.CheckAccess(ctx, req.UserId, req.CourseId)
	if err != nil {
		return nil, grpctransport.EncodeError(err)
	}
	return &purchasepb.CheckAccessResponse{HasAccess: hasAccess}, nil
}

func toProto(p model.Purchase) *purchasepb.Purchase {
	txID := ""
	if p.TransactionID != nil {
		txID = *p.TransactionID
	}
	return &purchasepb.Purchase{
		Id:            p.ID,
		UserId:        p.UserID,
		CourseId:      p.CourseID,
		Amount:        p.Amount,
		Status:        string(p.Status),
		TransactionId: txID,
	}
}
