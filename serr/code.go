package serr

type ErrorCode string

const (
	ErrCodeNotFound           ErrorCode = "NOT_FOUND"
	ErrCodeAlreadyExists      ErrorCode = "ALREADY_EXISTS"
	ErrCodePermissionDenied   ErrorCode = "PERMISSION_DENIED"
	ErrCodeResourceExhausted  ErrorCode = "RESOURCE_EXHAUSTED"
	ErrCodeInvalidArgument    ErrorCode = "INVALID_ARGUMENT"
	ErrCodeFailedPrecondition ErrorCode = "FAILED_PRECONDITION"
	ErrCodeDeadlineExceeded   ErrorCode = "DEADLINE_EXCEEDED"
	ErrCodeOutOfRange         ErrorCode = "OUT_OF_RANGE"
	ErrCodeUnimplemented      ErrorCode = "UNIMPLEMENTED"
	ErrCodeInternal           ErrorCode = "INTERNAL"
)
