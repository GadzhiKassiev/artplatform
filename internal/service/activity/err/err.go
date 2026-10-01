package err

import "artplatform/backend/internal/errs"

const (
	CodeStatsNotFound errs.Code = "STATS_NOT_FOUND"
)

var (
	ErrStatsNotFound = errs.New(CodeStatsNotFound, "course stats not found")
)
