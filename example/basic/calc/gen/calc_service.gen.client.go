// Code generated by Abide - DO NOT EDIT.
//
//	Timestamp: Mon, 03 Jul 2023 14:45:08 EDT
//	Source:    ../basic/calc/calc_service.go
//	Generator: https://github.com/monadicstack/abide
package calc

import (
	"context"

	"github.com/monadicstack/abide/example/basic/calc"
	"github.com/monadicstack/abide/fail"
	"github.com/monadicstack/abide/services/clients"
)

// CalculatorServiceClient creates an RPC client that conforms to the CalculatorService interface, but delegates
// work to remote instances. You must supply the base address of the remote service gateway instance or
// the load balancer for that service.
//
// CalculatorService provides the ability to perform basic arithmetic on two numbers.
func CalculatorServiceClient(address string, options ...clients.ClientOption) calc.CalculatorService {
	serviceClient := clients.NewClient("CalculatorService", address, options...)
	return &calculatorServiceClient{Client: serviceClient}
}

// calculatorServiceClient manages all interaction w/ a remote CalculatorService instance by letting you invoke functions
// on this instance as if you were doing it locally (hence... RPC client). Use the CalculatorServiceClient constructor
// function to actually get an instance of this client.
type calculatorServiceClient struct {
	clients.Client
}

// Add calculates and returns the sum of two numbers.
func (client *calculatorServiceClient) Add(ctx context.Context, request *calc.AddRequest) (*calc.AddResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &calc.AddResponse{}
	err := client.Invoke(ctx, "GET", "/add/{A}/{B}", request, response)
	return response, err

}

// Double multiplies the value by 2
func (client *calculatorServiceClient) Double(ctx context.Context, request *calc.DoubleRequest) (*calc.DoubleResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &calc.DoubleResponse{}
	err := client.Invoke(ctx, "POST", "/double/{Value}", request, response)
	return response, err

}

// Mul calculates and returns the product of two numbers.
func (client *calculatorServiceClient) Mul(ctx context.Context, request *calc.MulRequest) (*calc.MulResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &calc.MulResponse{}
	err := client.Invoke(ctx, "GET", "/multiply/{A}/{B}", request, response)
	return response, err

}

// Sub calculates and returns the difference between two numbers.
func (client *calculatorServiceClient) Sub(ctx context.Context, request *calc.SubRequest) (*calc.SubResponse, error) {

	if ctx == nil {
		return nil, fail.Unexpected("precondition failed: nil context")
	}
	if request == nil {
		return nil, fail.Unexpected("precondition failed: nil request")
	}

	response := &calc.SubResponse{}
	err := client.Invoke(ctx, "GET", "/sub/{A}/{B}", request, response)
	return response, err

}
