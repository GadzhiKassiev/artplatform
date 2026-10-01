package err

import "artplatform/backend/internal/errs"

const (
	CodeMediaNotFound   errs.Code = "MEDIA_NOT_FOUND"
	CodeUnsupportedType errs.Code = "UNSUPPORTED_TYPE"
	CodeStorageError    errs.Code = "STORAGE_ERROR"
	CodeAccessDenied    errs.Code = "ACCESS_DENIED"
	CodePreviewNotReady errs.Code = "PREVIEW_NOT_READY"
)

var (
	ErrMediaNotFound   = errs.New(CodeMediaNotFound, "media file not found")
	ErrUnsupportedType = errs.New(CodeUnsupportedType, "unsupported content type")
	ErrStorageError    = errs.New(CodeStorageError, "storage error")
	ErrAccessDenied    = errs.New(CodeAccessDenied, "access denied")
	ErrPreviewNotReady = errs.New(CodePreviewNotReady, "preview is not ready yet")
)
