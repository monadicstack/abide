package services

import (
	"context"
)

// MiddlewareFunc is a function that can be used to decorate a service method/endpoint's handler.
type MiddlewareFunc func(ctx context.Context, req any, next HandlerFunc) (any, error)

// MiddlewareFuncs is an ordered pipeline of operations that must occur before invoking
// a service method/endpoint's handler.
type MiddlewareFuncs []MiddlewareFunc

// Then creates a single handler function that executes every operation in the middleware
// pipeline and terminates with the supplied handler.
func (funcs MiddlewareFuncs) Then(handler HandlerFunc) HandlerFunc {
	for i := len(funcs) - 1; i >= 0; i-- {
		mw := funcs[i]
		next := handler
		handler = func(ctx context.Context, req any) (any, error) {
			return mw(ctx, req, next)
		}
	}
	return handler
}

// Append creates a new middleware function pipeline that runs the original handlers
// and then the additional ones specified by 'mw'.
func (funcs MiddlewareFuncs) Append(mw ...MiddlewareFunc) MiddlewareFuncs {
	return append(funcs, mw...)
}

// recoverMiddleware gets added as our outermost middleware to ensure that any accidental panic()
// calls at any level are gracefully caught without killing our server/process.
func recoverMiddleware() MiddlewareFunc {
	return func(ctx context.Context, req any, next HandlerFunc) (response any, err error) {
		defer func() {
			if recovery := recover(); recovery != nil {
				err, _ = recovery.(error)
			}
		}()
		return next(ctx, req)
	}
}
