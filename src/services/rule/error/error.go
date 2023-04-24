package error

import "insomnia/src/pkg/middleware_custom/grpc"

// define grpc business error code, the caller should handle error code

var (
	ParamError      = grpc.New(10001, "param error")
	PermissionError = grpc.New(10002, "permission forbidden")
	ConcurrentError = grpc.New(10003, "concurrent operation not allowed")
	RuleError       = grpc.New(10004, "linting rule content error")
)
