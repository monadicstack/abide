package apis

import (
	"net/http"

	"github.com/monadicstack/abide/codec"
	"github.com/monadicstack/abide/metadata"
	"github.com/monadicstack/abide/services"
)

// HTTPMiddlewareFunc is a function that can intercept HTTP requests going through your API
// gateway, allowing you to directly manipulate the HTTP request/response as needed.
type HTTPMiddlewareFunc func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc)

// HTTPMiddlewareFuncs defines an ordered pipeline of middleware functions you want an endpoint
// in your API gateway to perform.
type HTTPMiddlewareFuncs []HTTPMiddlewareFunc

// Append adds the given function(s) to the END of the receiver's current pipeline.
//
//	middlewareFuncs := HTTPMiddlewareFuncs{A, B, C}
//	middlewareFuncs = middlewareFuncs.Append(D, E, F)
//	// middlewareFuncs ==> {A, B, C, D, E, F}
func (funcs HTTPMiddlewareFuncs) Append(functions ...HTTPMiddlewareFunc) HTTPMiddlewareFuncs {
	return append(funcs, functions...)
}

// Then converts an entire middleware chain/pipeline into a single handler function that can be
// registered with a server/router. The 'handler' parameter goes at the very end of the chain; the
// idea being that you perform all of your middleware tasks and 'handler' is the "real work" you
// wanted to accomplish all along.
func (funcs HTTPMiddlewareFuncs) Then(handler http.HandlerFunc) http.HandlerFunc {
	for i := len(funcs) - 1; i >= 0; i-- {
		mw := funcs[i]
		next := handler
		handler = func(w http.ResponseWriter, req *http.Request) {
			mw(w, req, next)
		}
	}
	return handler
}

// recoverFromPanic automatically recovers from a panic thrown by your handler so that if you nil-pointer
// or something else unexpected, we'll safely just return a 500-style error.
func recoverFromPanic(encoder codec.Encoder) HTTPMiddlewareFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Content-Type", encoder.ContentType())
				w.WriteHeader(http.StatusInternalServerError)
				_ = encoder.Encode(w, err)
			}
		}()
		next(w, req)
	}
}

// restoreMetadata looks for the X-RPC-Metadata header, decodes it, and places the appropriate
// metadata values back onto the request context so the rest of the operation already has access
// to them. This is how Service B automatically has access to the same auth/values/etc. when
// called from Service A.
func restoreMetadata() HTTPMiddlewareFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		encodedMetadata := metadata.EncodedBytes(req.Header.Get("X-RPC-Metadata"))
		ctx := metadata.Decode(req.Context(), encodedMetadata)
		next(w, req.WithContext(ctx))
	}
}

// restoreMetadataEndpoint simply adds the routing metadata, so that you can determine
// which service operation you're calling from any of your general purpose metadata.
func restoreMetadataEndpoint(endpoint services.Endpoint, route services.EndpointRoute) HTTPMiddlewareFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		// Even if the context has this info, it's probably from another service call whose
		// route is different from this one. The route needs to be for THIS service call.
		ctx := metadata.WithRoute(req.Context(), metadata.EndpointRoute{
			ServiceName: endpoint.ServiceName,
			Name:        endpoint.Name,
			Type:        route.GatewayType.String(),
			Method:      route.Method,
			Path:        route.Path,
		})
		next(w, req.WithContext(ctx))
	}
}

// restoreAuthorization applies the Authorization HTTP header to your context metadata.
func restoreAuthorization() HTTPMiddlewareFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		// It is possible to get authorization from the metadata attributes header as well
		// as the standard Authorization header. We will use the metadata value if there's
		// not one in the standard header.
		//
		// This allows you to authorize a request using one set of credentials and then that
		// call can invoke a call on another service using a completely different set of
		// credentials. If we used the metadata attrs as the authoritative value, you'd be
		// stuck using one set of credentials for everything this request may require which
		// is not what we want.
		ctx := req.Context()
		if auth := req.Header.Get("Authorization"); auth != "" {
			ctx = metadata.WithAuthorization(req.Context(), auth)
		}
		next(w, req.WithContext(ctx))
	}
}

// restoreTraceID ensures that this request ALWAYS has a unique request/trace id for use in
// your logging/observability code. It will restore the value provided by some downstream
// service if present; otherwise it will generate a unique-enough value for you.
func restoreTraceID() HTTPMiddlewareFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		// We're already carrying over a request id from a previous call in this chain.
		if traceID := metadata.TraceID(req.Context()); traceID != "" {
			next(w, req)
			return
		}

		// This is probably the primordial service call, so use the HTTP header value if
		// there is one - otherwise, generate one for us to use. All requests should have one.
		switch traceID := req.Header.Get("X-Request-ID"); traceID {
		case "":
			ctx := metadata.WithTraceID(req.Context(), metadata.NewTraceID())
			next(w, req.WithContext(ctx))
		default:
			ctx := metadata.WithTraceID(req.Context(), traceID)
			next(w, req.WithContext(ctx))
		}
	}
}

// restoreMetadataHeaders places the HTTP header map into the request metadata. This way
// you can tweak service behavior based on things like Cache-Control or things like that.
func restoreMetadataHeaders() HTTPMiddlewareFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		ctx := metadata.WithRequestHeaders(req.Context(), req.Header)
		next(w, req.WithContext(ctx))
	}
}
