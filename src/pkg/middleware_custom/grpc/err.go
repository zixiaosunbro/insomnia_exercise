package grpc

import (
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorCode is type for error code.
type ErrorCode uint32

const (
	CodeCodeUnknownError ErrorCode = 100
)

type CooError struct {
	Code ErrorCode
	Msg  string
}

func (n CooError) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", n.Code, n.Msg)
}

func SoaErr(t CooError) error {
	return status.Errorf(codes.Code(t.Code), t.Msg)
}

func New(code int, msg string) CooError {
	return CooError{Code: ErrorCode(code), Msg: msg}
}
