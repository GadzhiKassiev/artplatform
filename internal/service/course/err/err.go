package err

import "artplatform/backend/internal/errs"

const (
	CodeNotOwner         errs.Code = "NOT_OWNER"
	CodeAlreadyPublished errs.Code = "ALREADY_PUBLISHED"
	CodeInvalidPrice     errs.Code = "INVALID_PRICE"
	CodeEmptyTitle       errs.Code = "EMPTY_TITLE"
)

var (
	ErrNotOwner         = errs.New(CodeNotOwner, "only author can modify course")
	ErrAlreadyPublished = errs.New(CodeAlreadyPublished, "course already published")
	ErrInvalidPrice     = errs.New(CodeInvalidPrice, "price must be >= 0")
	ErrEmptyTitle       = errs.New(CodeEmptyTitle, "title must not be empty")
)
