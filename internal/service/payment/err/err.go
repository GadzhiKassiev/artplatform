package err

import "artplatform/backend/internal/errs"

const (
	CodeTransactionNotFound errs.Code = "TRANSACTION_NOT_FOUND"
	CodePaymentFailed       errs.Code = "PAYMENT_FAILED"
)

var (
	ErrTransactionNotFound = errs.New(CodeTransactionNotFound, "transaction not found")
	ErrPaymentFailed       = errs.New(CodePaymentFailed, "payment processing failed")
)
