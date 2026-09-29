package err

import "artplatform/backend/internal/errs"

const (
	CodeEmailTaken         errs.Code = "EMAIL_TAKEN"
	CodeInvalidCredentials errs.Code = "INVALID_CREDENTIALS"
	CodeInvalidRole        errs.Code = "INVALID_ROLE"
)

var (
	ErrEmailTaken         = errs.New(CodeEmailTaken, "email already taken")
	ErrInvalidCredentials = errs.New(CodeInvalidCredentials, "invalid credentials")
	ErrInvalidRole        = errs.New(CodeInvalidRole, "invalid role")
)
