package graphql

import (
	"artplatform/backend/api/graphql/generated/model"
	purchasepb "artplatform/backend/proto/purchase"
)

func toPurchaseModel(p *purchasepb.Purchase) *model.Purchase {
	if p == nil {
		return nil
	}
	return &model.Purchase{
		ID:            p.Id,
		UserID:        p.UserId,
		CourseID:      p.CourseId,
		Amount:        p.Amount,
		Status:        p.Status,
		TransactionID: strPtr(p.TransactionId),
	}
}
