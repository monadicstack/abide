//go:build integration

package services_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/monadicstack/abide/internal/quiet"
	"github.com/monadicstack/abide/internal/testext"
	gen "github.com/monadicstack/abide/internal/testext/gen"
	"github.com/monadicstack/abide/services"
	"github.com/monadicstack/abide/services/gateways/apis"
	"github.com/monadicstack/abide/services/gateways/events"
	"github.com/stretchr/testify/suite"
)

func TestServerSuite(t *testing.T) {
	suite.Run(t, &ServerSuite{addresses: testext.NewFreeAddress("localhost", 20000)})
}

type ServerSuite struct {
	suite.Suite
	addresses testext.FreeAddress
}

func (suite *ServerSuite) start() (*services.Server, *testext.Sequence, func()) {
	// Grab a fresh address, so we can parallelize our tests.
	address := suite.addresses.Next()

	// Capture invocations across both services in one timeline.
	sequence := &testext.Sequence{}

	sampleService := testext.SampleServiceHandler{Sequence: sequence}
	otherService := testext.OtherServiceHandler{
		Sequence:      sequence,
		SampleService: gen.SampleServiceClient(address),
	}

	server := services.NewServer(
		services.Listen(apis.NewGateway(address)),
		services.Listen(events.NewGateway()),
		services.Register(gen.SampleServiceServer(sampleService)),
		services.Register(gen.OtherServiceServer(otherService)),
	)
	go server.Run()

	// Kinda crappy, but we need some time to make sure the server is up. Sometimes
	// this goes so fast that the test case fires before the server is fully running.
	// As a result the cases fail because the server's not running... duh.
	time.Sleep(25 * time.Millisecond)

	return server, sequence, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}
}

func (suite *ServerSuite) responseText(res any) string {
	if sampleRes, ok := res.(*testext.SampleResponse); ok {
		return sampleRes.Text
	}
	if otherRes, ok := res.(*testext.OtherResponse); ok {
		return otherRes.Text
	}
	return "<invalid response type>"
}

func (suite *ServerSuite) responseStream(res any) *services.StreamResponse {
	if res == nil {
		return nil
	}
	if sampleRes, ok := res.(*testext.SampleDownloadResponse); ok {
		return &sampleRes.StreamResponse
	}
	if sampleRes, ok := res.(*testext.SampleRedirectResponse); ok {
		return &sampleRes.StreamResponse
	}
	return nil
}

func (suite *ServerSuite) streamContent(stream *services.StreamResponse) string {
	if stream == nil {
		return ""
	}
	content := stream.Content()
	if content == nil {
		return ""
	}

	defer quiet.Close(content)
	data, err := io.ReadAll(content)
	if err != nil {
		return ""
	}
	return string(data)
}

func (suite *ServerSuite) assertInvoked(calls *testext.Sequence, expected []string) {
	time.Sleep(25 * time.Millisecond)
	suite.ElementsMatch(calls.Values(), expected)
}

func (suite *ServerSuite) TestBasicExecution() {
	server, _, shutdown := suite.start()
	defer shutdown()

	res, err := server.Invoke(context.Background(), "SampleService", "Defaults", &testext.SampleRequest{
		Text: "Hello",
	})
	suite.Require().NoError(err)
	suite.Equal("Defaults:Hello", suite.responseText(res))
}

func (suite *ServerSuite) TestStreamedResponse() {
	server, _, shutdown := suite.start()
	defer shutdown()

	res, err := server.Invoke(context.Background(), "SampleService", "Download", &testext.SampleDownloadRequest{
		Format: "text/plain",
	})
	suite.Require().NoError(err)
	stream := suite.responseStream(res)
	suite.Equal("Donny, you're out of your element!", suite.streamContent(stream))
	suite.Equal("text/plain", stream.ContentType())
	suite.Equal(34, stream.ContentLength())
}

func (suite *ServerSuite) TestResumableStreamedResponse() {
	server, _, shutdown := suite.start()
	defer shutdown()

	res, err := server.Invoke(context.Background(), "SampleService", "DownloadResumable", &testext.SampleDownloadRequest{
		Format: "text/plain",
	})
	suite.Require().NoError(err)
	stream := suite.responseStream(res)
	suite.Equal("<h1>The Dude Abides</h1>", suite.streamContent(stream))
	suite.Equal("text/html", stream.ContentType())

	start, end, size := stream.ContentRange()
	suite.Equal(50, start)
	suite.Equal(74, end)
	suite.Equal(1024, size)
}

// Ensure that the service still manages an endpoint even if it's not exposed by any gateway.
func (suite *ServerSuite) TestOmittedEndpoint() {
	server, _, shutdown := suite.start()
	defer shutdown()

	res, err := server.Invoke(context.Background(), "SampleService", "OmitMe", &testext.SampleRequest{
		Text: "Hello",
	})
	suite.Require().NoError(err)
	suite.Equal("Doesn't matter...", suite.responseText(res))
}

// Make sure that regardless of how endpoints are invoked, they trigger gateway
// specific middleware such as event publishing.
func (suite *ServerSuite) TestEvents() {
	server, calls, shutdown := suite.start()
	defer shutdown()

	calls.Reset()
	res, err := server.Invoke(context.Background(), "SampleService", "TriggerLowerCase", &testext.SampleRequest{Text: "Abide"})
	suite.Require().NoError(err)
	suite.Equal("abide", suite.responseText(res))
	suite.assertInvoked(calls, []string{
		"TriggerLowerCase:Abide",
		"ListenerB:abide", // event should receive the result of TriggerLowerCase
	})

	calls.Reset()
	res, err = server.Invoke(context.Background(), "SampleService", "TriggerError", &testext.SampleRequest{Text: "Abide"})
	suite.Require().Error(err)
	suite.assertInvoked(calls, []string{})
}

// Ensure that ServiceA is able to listen to events from ServiceB and that ServiceB
// can listen to events from ServiceA as well. As a side effect, this one also
// makes sure that event triggers and cascade and cause others to trigger.
func (suite *ServerSuite) TestEvents_bidirectional() {
	server, calls, shutdown := suite.start()
	defer shutdown()

	calls.Reset()
	sampleRes, err := server.Invoke(context.Background(), "SampleService", "TriggerUpperCase", &testext.SampleRequest{Text: "Abide"})
	suite.Require().NoError(err)
	suite.Equal("ABIDE", suite.responseText(sampleRes))
	suite.assertInvoked(calls, []string{
		"TriggerUpperCase:Abide",
		"ListenerA:ABIDE",
		"ListenerB:ABIDE",
		"ListenerB:ListenerA:ABIDE", // Cascade when ListenerB fired after ListenerA fired
		"ListenWell:ABIDE",          // OtherService also gets SampleService.TriggerUpperCase events
	})

	calls.Reset()
	otherRes, err := server.Invoke(context.Background(), "OtherService", "SpaceOut", &testext.OtherRequest{Text: "Abide"})
	suite.Require().NoError(err)
	suite.Equal("A b i d e", suite.responseText(otherRes))
	suite.assertInvoked(calls, []string{
		"SpaceOut:Abide",       // The original call we invoked.
		"ListenerB:A b i d e",  // The subscriber in the SampleService
		"ListenWell:A b i d e", // The subscriber in the OtherService
	})
}

// Ensure that one service can invoke functions on another and that those functions can cause
// events to fire and be handled on multiple services.
func (suite *ServerSuite) TestRPCWithEvents() {
	server, calls, shutdown := suite.start()
	defer shutdown()

	calls.Reset()
	res, err := server.Invoke(context.Background(), "OtherService", "RPCExample", &testext.OtherRequest{Text: "Abide"})
	suite.Require().NoError(err)
	suite.Equal("ABIDE", suite.responseText(res))
	suite.assertInvoked(calls, []string{
		// Our initial call a few lines above.
		"RPCExample:Abide",

		// The RPC call that the other service made to the sample service.
		"TriggerUpperCase:Abide",

		// Since RPCExample explicitly calls TriggerUpperCase, when *that* call finishes, it
		// should trigger events that are handled by methods on both services.
		"ListenerA:ABIDE",
		"ListenerB:ABIDE",
		"ListenerB:ListenerA:ABIDE", // Cascade when ListenerB fired after ListenerA fired
		"ListenWell:ABIDE",
	})
}
