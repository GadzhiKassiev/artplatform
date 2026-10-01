package err

import "artplatform/backend/internal/errs"

const (
	CodePurchaseNotFound errs.Code = "PURCHASE_NOT_FOUND"
	CodeAlreadyPurchased errs.Code = "ALREADY_PURCHASED"
	CodePaymentFailed    errs.Code = "PAYMENT_FAILED"
)

var (
	ErrPurchaseNotFound = errs.New(CodePurchaseNotFound, "purchase not found")
	ErrAlreadyPurchased = errs.New(CodeAlreadyPurchased, "course already purchased")
	ErrPaymentFailed    = errs.New(CodePaymentFailed, "payment failed")
)
