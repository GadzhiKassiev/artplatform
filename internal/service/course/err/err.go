package err

import "artplatform/backend/internal/errs"

const (
	CodeNotOwner         errs.Code = "NOT_OWNER"
	CodeAlreadyPublished errs.Code = "ALREADY_PUBLISHED"
)

var (
	ErrNotOwner         = errs.New(CodeNotOwner, "only author can modify course")
	ErrAlreadyPublished = errs.New(CodeAlreadyPublished, "course already published")
)
