// Code generated by Abide - DO NOT EDIT.
//
//	Timestamp: Thu, 29 Jun 2023 09:39:40 EDT
//	Source:    sample_service.go
//	Generator: https://github.com/monadicstack/abide
package testext

import (
	"context"

	"github.com/monadicstack/abide/fail"
	"github.com/monadicstack/abide/internal/testext"
	"github.com/monadicstack/abide/services/clients"
)

// SampleServiceClient creates an RPC client that conforms to the SampleService interface, but delegates
// work to remote instances. You must supply the base address of the remote service gateway instance or
// the load balancer for that service.
//
// SampleService is a mix of different options, parameter setups, and responses so that we can
// run integration tests using our code-generated clients. Each method is nothing special, but
// they each do something a little differently than the rest to flex different parts of the framework.
func SampleServiceClient(address string, options ...clients.ClientOption) testext.SampleService {
	serviceClient := clients.NewClient("SampleService", address, options...)
	return &sampleServiceClient{Client: serviceClient}
}

// sampleServiceClient manages all interaction w/ a remote SampleService instance by letting you invoke functions
// on this instance as if you were doing it locally (hence... RPC client). Use the SampleServiceClient constructor
// function to actually get an instance of this client.
type sampleServiceClient struct {
	clients.Client
}

// Authorization regurgitates the "Authorization" metadata/header.
func (client *sampleServiceClient) Authorization(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "POST", "/v2/SampleService.Authorization", request, response)
	return response, err

}

// ComplexValues flexes our ability to encode/decode non-flat structs.
func (client *sampleServiceClient) ComplexValues(ctx context.Context, request *testext.SampleComplexRequest) (*testext.SampleComplexResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleComplexResponse{}
	err := client.Invoke(ctx, "POST", "/v2/SampleService.ComplexValues", request, response)
	return response, err

}

// ComplexValuesPath flexes our ability to encode/decode non-flat structs while
// specifying them via path and query string.
func (client *sampleServiceClient) ComplexValuesPath(ctx context.Context, request *testext.SampleComplexRequest) (*testext.SampleComplexResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleComplexResponse{}
	err := client.Invoke(ctx, "GET", "/v2/complex/values/{InUser.ID}/{InUser.Name}/woot", request, response)
	return response, err

}

// CustomRoute performs a service operation where you override default behavior
// by providing routing-related Doc Options.
func (client *sampleServiceClient) CustomRoute(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "GET", "/v2/custom/route/1/{ID}/{Text}", request, response)
	return response, err

}

// CustomRouteBody performs a service operation where you override default behavior
// by providing routing-related Doc Options, but rely on body encoding rather than path.
func (client *sampleServiceClient) CustomRouteBody(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "PUT", "/v2/custom/route/3/{ID}", request, response)
	return response, err

}

// CustomRouteQuery performs a service operation where you override default behavior
// by providing routing-related Doc Options. The input data relies on the path
func (client *sampleServiceClient) CustomRouteQuery(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "GET", "/v2/custom/route/2/{ID}", request, response)
	return response, err

}

// Defaults simply utilizes all of the framework's default behaviors.
func (client *sampleServiceClient) Defaults(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "POST", "/v2/SampleService.Defaults", request, response)
	return response, err

}

// Download results in a raw stream of data rather than relying on auto-encoding
// the response value.
func (client *sampleServiceClient) Download(ctx context.Context, request *testext.SampleDownloadRequest) (*testext.SampleDownloadResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleDownloadResponse{}
	err := client.Invoke(ctx, "GET", "/v2/download", request, response)
	return response, err

}

// DownloadResumable results in a raw stream of data rather than relying on auto-encoding
// the response value. The stream includes Content-Range info as though you could resume
// your stream/download progress later.
func (client *sampleServiceClient) DownloadResumable(ctx context.Context, request *testext.SampleDownloadRequest) (*testext.SampleDownloadResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleDownloadResponse{}
	err := client.Invoke(ctx, "GET", "/v2/download/resumable", request, response)
	return response, err

}

// Fail4XX always returns a non-nil 400-series error.
func (client *sampleServiceClient) Fail4XX(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "POST", "/v2/SampleService.Fail4XX", request, response)
	return response, err

}

// Fail5XX always returns a non-nil 500-series error.
func (client *sampleServiceClient) Fail5XX(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "POST", "/v2/SampleService.Fail5XX", request, response)
	return response, err

}

// ListenerA fires on only one of the triggers.
func (client *sampleServiceClient) ListenerA(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "GET", "/v2/ListenerA/Woot", request, response)
	return response, err

}

// ListenerB fires on multiple triggers... including another event-based endpoint. We also
// listen for the TriggerFailure event which should never fire properly.
func (client *sampleServiceClient) ListenerB(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	// Not exposed, so don't bother with a round trip to the server just to get a "not found" error anyway.
	return nil, fail.NotImplemented("ListenerB is not supported in the API gateway")

}

// OmitMe exists in the service, but should be excluded from the public API.
func (client *sampleServiceClient) OmitMe(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	// Not exposed, so don't bother with a round trip to the server just to get a "not found" error anyway.
	return nil, fail.NotImplemented("OmitMe is not supported in the API gateway")

}

// Redirect results in a 307-style redirect to the Download endpoint.
func (client *sampleServiceClient) Redirect(ctx context.Context, request *testext.SampleRedirectRequest) (*testext.SampleRedirectResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleRedirectResponse{}
	err := client.Invoke(ctx, "GET", "/v2/redirect", request, response)
	return response, err

}

// Sleep successfully responds, but it will sleep for 5 seconds before doing so. Use this
// for test cases where you want to try out timeouts.
func (client *sampleServiceClient) Sleep(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "POST", "/v2/SampleService.Sleep", request, response)
	return response, err

}

func (client *sampleServiceClient) TriggerFailure(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "POST", "/v2/SampleService.TriggerFailure", request, response)
	return response, err

}

func (client *sampleServiceClient) TriggerLowerCase(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "POST", "/v2/SampleService.TriggerLowerCase", request, response)
	return response, err

}

// TriggerUpperCase ensures that events still fire as "SampleService.TriggerUpperCase" even though
// we are going to set a different HTTP path.
func (client *sampleServiceClient) TriggerUpperCase(ctx context.Context, request *testext.SampleRequest) (*testext.SampleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.SampleResponse{}
	err := client.Invoke(ctx, "GET", "/v2/Upper/Case/WootyAndTheBlowfish", request, response)
	return response, err

}
