package grpc

import (
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"artplatform/backend/internal/errs"
)

func EncodeError(err error) error {
	if err == nil {
		return nil
	}

	var domainErr *errs.Error
	if errors.As(err, &domainErr) {
		return status.Error(codes.InvalidArgument, string(domainErr.Code)+"|"+domainErr.Message)
	}

	return status.Error(codes.Internal, "internal error")
}

func DecodeError(err error) error {
	if err == nil {
		return nil
	}

	if st, ok := status.FromError(err); ok {
		if de := parsePattern(st.Message()); de != nil {
			return de
		}
	}

	if de := parsePattern(err.Error()); de != nil {
		return de
	}

	return err
}

func parsePattern(s string) *errs.Error {
	idx := strings.Index(s, "|")
	if idx <= 0 {
		return nil
	}

	codeStart := idx
	for codeStart > 0 {
		c := s[codeStart-1]
		if c == ' ' || c == '=' || c == ':' {
			break
		}
		codeStart--
	}

	codeStr := s[codeStart:idx]
	message := s[idx+1:]

	if !isLikelyCode(codeStr) {
		return nil
	}

	return errs.New(errs.Code(codeStr), message)
}

func isLikelyCode(s string) bool {
	if len(s) == 0 || len(s) > 50 {
		return false
	}
	for _, c := range s {
		if !((c >= 'A' && c <= 'Z') || c == '_') {
			return false
		}
	}
	return true
}
