package serr

import "context"

type Options struct {
	err  error
	code ErrorCode
	ctx  context.Context
}

type Option func(*Options)

func WithError(err error) Option {
	return func(o *Options) {
		o.err = err
	}
}

func WithErrorCode(code ErrorCode) Option {
	return func(o *Options) {
		o.code = code
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.ctx = ctx
	}
}
