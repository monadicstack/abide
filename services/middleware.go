package services

import (
	"context"

	"github.com/monadicstack/abide/internal/naming"
	"github.com/monadicstack/abide/internal/reflection"
	"github.com/monadicstack/abide/internal/slices"
	"github.com/monadicstack/abide/metadata"
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

// rolesMiddleware takes the raw doc option roles list such as ["admin.write", "group.{ID}.write"] and populates
// the path variables w/ runtime values, so you end up with a roles list like ["admin.write", "group.123.write"]. For
// any path variables that can't be properly mapped to a runtime value, those will end up blank (e.g. "group..write").
func rolesMiddleware(endpoint Endpoint) MiddlewareFunc {
	return func(ctx context.Context, req any, next HandlerFunc) (any, error) {
		populateRole := func(role string) string {
			return naming.ResolvePath(role, '.', func(variable string) string {
				var runtimeValue string
				reflection.ToBindingValue(req, variable, &runtimeValue)
				return runtimeValue
			})
		}

		route := metadata.Route(ctx)
		route.Roles = slices.Map(endpoint.Roles, populateRole)
		return next(metadata.WithRoute(ctx, route), req)
	}
}
