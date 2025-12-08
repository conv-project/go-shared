package serr

import (
	"context"
	"fmt"
)

type AppError struct {
	Code    ErrorCode
	Message string // безопасное описание
	Err     error  // внутренняя причина (wrapped)
	Context context.Context
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

func New(message string, opts ...Option) *AppError {
	o := Options{
		code: ErrCodeInternal,
		ctx:  context.Background(),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &AppError{
		Code:    o.code,
		Message: message,
		Err:     o.err,
		Context: o.ctx,
	}
}
