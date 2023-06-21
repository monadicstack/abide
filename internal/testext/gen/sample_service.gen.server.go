// Code generated by Boson - DO NOT EDIT.
//
//	Timestamp: Wed, 21 Jun 2023 16:34:43 EDT
//	Source:    sample_service.go
//	Generator: https://github.com/monadicstack/abide
package testext

import (
	"context"

	"github.com/monadicstack/abide/fail"
	"github.com/monadicstack/abide/internal/testext"
	"github.com/monadicstack/abide/services"
)

// SampleServiceServer accepts your "real" SampleService instance (the thing that really does
// the work), and returns a set of endpoint routes which allow this service to be consumed
// via the gateways/listeners you configure in main().
//
//	// Example
//	serviceHandler := testext.SampleServiceHandler{ /* set up to your liking */ }
//	server := services.New(
//		services.Listen(apis.NewGateway()),
//		services.Register(testextgen.SampleServiceServer(serviceHandler)),
//	)
//	server.Listen()
//
// From there, you can add middleware, event sourcing support and more. Look at the abide
// documentation for more details/examples on how to make your service production ready.
func SampleServiceServer(handler testext.SampleService, middleware ...services.MiddlewareFunc) *services.Service {
	middlewareFuncs := services.MiddlewareFuncs(middleware)

	return &services.Service{
		Name:    "SampleService",
		Version: "0.0.1",
		Handler: handler,
		Endpoints: []services.Endpoint{

			{
				ServiceName: "SampleService",
				Name:        "Authorization",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Authorization(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Authorization",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "ComplexValues",
				NewInput:    func() services.StructPointer { return &testext.SampleComplexRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleComplexRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.ComplexValues(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.ComplexValues",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "ComplexValuesPath",
				NewInput:    func() services.StructPointer { return &testext.SampleComplexRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleComplexRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.ComplexValuesPath(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/complex/values/:InUser.ID/:InUser.Name/woot",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "CustomRoute",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.CustomRoute(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/custom/route/1/:ID/:Text",
						Status:      202,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "CustomRouteBody",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.CustomRouteBody(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "PUT",
						Path:        "/v2/custom/route/3/:ID",
						Status:      201,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "CustomRouteQuery",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.CustomRouteQuery(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/custom/route/2/:ID",
						Status:      202,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Defaults",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Defaults(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Defaults",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Download",
				NewInput:    func() services.StructPointer { return &testext.SampleDownloadRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleDownloadRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Download(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/download",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "DownloadResumable",
				NewInput:    func() services.StructPointer { return &testext.SampleDownloadRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleDownloadRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.DownloadResumable(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/download/resumable",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Fail4XX",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Fail4XX(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Fail4XX",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Fail5XX",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Fail5XX(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Fail5XX",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "ListenerA",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.ListenerA(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.ListenerA",
						Status:      200,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.TriggerUpperCase",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "ListenerB",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.ListenerB(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.TriggerUpperCase",
						Status:      0,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.TriggerLowerCase",
						Status:      0,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.TriggerFailure",
						Status:      0,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "SampleService.ListenerA",
						Status:      0,
					},

					{
						GatewayType: "EVENTS",
						Method:      "ON",
						Path:        "OtherService.SpaceOut",
						Status:      0,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "OmitMe",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.OmitMe(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{},
			},

			{
				ServiceName: "SampleService",
				Name:        "Redirect",
				NewInput:    func() services.StructPointer { return &testext.SampleRedirectRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRedirectRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Redirect(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "GET",
						Path:        "/v2/redirect",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "Sleep",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.Sleep(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.Sleep",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "TriggerFailure",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.TriggerFailure(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.TriggerFailure",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "TriggerLowerCase",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.TriggerLowerCase(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.TriggerLowerCase",
						Status:      200,
					},
				},
			},

			{
				ServiceName: "SampleService",
				Name:        "TriggerUpperCase",
				NewInput:    func() services.StructPointer { return &testext.SampleRequest{} },
				Handler: middlewareFuncs.Then(func(ctx context.Context, req any) (any, error) {
					typedReq, ok := req.(*testext.SampleRequest)
					if !ok {
						return nil, fail.Unexpected("invalid request argument type")
					}
					return handler.TriggerUpperCase(ctx, typedReq)
				}),
				Routes: []services.EndpointRoute{

					{
						GatewayType: "API",
						Method:      "POST",
						Path:        "/v2/SampleService.TriggerUpperCase",
						Status:      200,
					},
				},
			},
		},
	}
}
