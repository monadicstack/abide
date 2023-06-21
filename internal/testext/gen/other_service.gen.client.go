// Code generated by Boson - DO NOT EDIT.
//
//	Timestamp: Wed, 21 Jun 2023 16:34:43 EDT
//	Source:    other_service.go
//	Generator: https://github.com/monadicstack/abide
package testext

import (
	"context"

	"github.com/monadicstack/abide/fail"
	"github.com/monadicstack/abide/internal/testext"
	"github.com/monadicstack/abide/services/clients"
)

// OtherServiceClient creates an RPC client that conforms to the OtherService interface, but delegates
// work to remote instances. You must supply the base address of the remote service gateway instance or
// the load balancer for that service.
//
// OtherService primarily exists to show that we can send event signals between services.
func OtherServiceClient(address string, options ...clients.ClientOption) testext.OtherService {
	serviceClient := clients.NewClient("OtherService", address, options...)
	return &otherServiceClient{Client: serviceClient}
}

// otherServiceClient manages all interaction w/ a remote OtherService instance by letting you invoke functions
// on this instance as if you were doing it locally (hence... RPC client). Use the OtherServiceClient constructor
// function to actually get an instance of this client.
type otherServiceClient struct {
	clients.Client
}

// ListenWell can listen for successful responses across multiple services.
func (client *otherServiceClient) ListenWell(ctx context.Context, request *testext.OtherRequest) (*testext.OtherResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.OtherResponse{}
	err := client.Invoke(ctx, "POST", "/OtherService.ListenWell", request, response)
	return response, err

}

// RPCExample invokes the TriggerUpperCase() function on the SampleService to get work done.
// This will make sure that we can do cross-service communication.
func (client *otherServiceClient) RPCExample(ctx context.Context, request *testext.OtherRequest) (*testext.OtherResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.OtherResponse{}
	err := client.Invoke(ctx, "POST", "/OtherService.RPCExample", request, response)
	return response, err

}

// SpaceOut takes your input text and puts spaces in between all the letters.
func (client *otherServiceClient) SpaceOut(ctx context.Context, request *testext.OtherRequest) (*testext.OtherResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &testext.OtherResponse{}
	err := client.Invoke(ctx, "POST", "/OtherService.SpaceOut", request, response)
	return response, err

}
