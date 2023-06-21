# Abide

Abide is a code generator and runtime library that helps you
write (micro) services/APIs that supports, both, RPC/HTTP
and Event-Driven invocation. It parses the
interfaces/structs/comments in your code service
code to generate all of the client, server, gateway, and pub-sub
communication code automatically.

You write business logic. Abide generates the annoying
copy/paste boilerplate needed to expose your service as an
HTTP API as well as Pub/Sub code to create event-driven
workflows across your services.

This is the spiritual successor to [Frodo](https://github.com/monadicstack/frodo).
Abide supports Frodo's RPC/HTTP related features, but it addresses
many shortcomings in the architecture/approach and adds
Event-Driven communication with almost no extra code on your part.

## Getting Started

```shell
go install github.com/monadicstack/abide@latest
```
This will fetch the `abide` code generation executable as well
as the runtime libraries that allow your services and clients to
communicate with each other.

## Basic Example

We're going to write a simple `CalculatorService` that
lets you either add or subtract two numbers.

### Step 1: Describe Your Service

Your first step is to write a .go file that just defines
the contract for your service; the interface as well as the
inputs/outputs.

```go
// calc/calculator_service.go

package calc

import (
    "context"
)

type CalculatorService interface {
    Add(context.Context, *AddRequest) (*AddResponse, error)
    Sub(context.Context, *SubRequest) (*SubResponse, error)
}

type AddRequest struct {
    A int
    B int
}

type AddResponse struct {
    Result int
}

type SubRequest struct {
    A int
    B int
}

type SubResponse struct {
    Result int
}
```

One important detail is that the interface name ends with
the suffix "Service". This tells Abide that this is an
actual service interface and not just some random abstraction
in your code.

At this point you haven't actually defined *how* this service gets
this work done; just which operations are available.

### Step 2: Implement Your Service Logic

We actually have enough for Abide to
generate your RPC/API/Event code already, but we'll hold off
for a moment. Abide frees you up to focus on building
features, so let's actually implement service - no networking,
no marshaling, no status stuff, no pub/sub - just business logic to make your
service behave properly.

```go
// calc/calculator_service_handler.go

package calc

import (
    "context"
)

type CalculatorServiceHandler struct {}

func (svc CalculatorServiceHandler) Add(ctx context.Context, req *AddRequest) (*AddResponse, error) {
    result := req.A + req.B
    return &AddResponse{Result: result}, nil
}

func (svc CalculatorServiceHandler) Sub(ctx context.Context, req *SubRequest) (*SubResponse, error) {
    result := req.A - req.B
    return &SubResponse{Result: result}, nil
}
```

### Step 3: Generate Your RPC Client and Server Code

At this point, you've just written the same code that you (hopefully)
would have written even if you weren't using Abide. Next,
we want to auto-generate two things:

* The "server" bits that allow an instance of your CalculatorService
  to listen for incoming requests from an either HTTP API or
  a published event. (We'll look at events later...)
* A "client" struct that communicates with that API to get work done.

Just run these two commands in a terminal:

```shell
# Feed it the service interface file, not the handler.
abide server calculator_service.go
abide client calculator_service.go
```

### Step 4: Run Your Calculator API

Let's fire up an HTTP server on port 9000 that makes your service
available for consumption.

```go
package main

import (
    "github.com/monadicstack/abide/services"
    "github.com/monadicstack/abide/services/gateways/apis"

    "github.com/your/project/calc"
    calcgen "github.com/your/project/calc/gen"
)

func main() {
    // Create your logic-only handler, then wrap it in service
    // communication bits that let it interact with the Abide runtime.
    calcHandler := calc.CalculatorServiceHandler{}
    calcService := calcgen.NewCalculatorService(calcHandler)
	
    // Fire up a server that will manage our service and listen to 
    // API calls on port 9000.
    server := runtime.New(calcService)
    server.Listen(apis.NewGateway(":9000"))
    server.Run()
}
```

Seriously. That's the whole program.

Compile and run it, and your service/API is now ready
to be consumed. We'll use the Go client we generated in just
a moment, but you can try this out right now by simply
using curl:

```shell
curl -d '{"A":5, "B":2}' http://localhost:9000/CalculatorService.Add
# {"Result":7}
curl -d '{"A":5, "B":2}' http://localhost:9000/CalculatorService.Sub
# {"Result":3}
```

### Step 5: Interact With Your Calculator Service

While you can use raw HTTP to communicate with the service,
let's use our auto-generated client to hide the gory
details of JSON marshaling, status code translation, and
other noise.

The client actually implements CalculatorService
just like the server/handler does. As a result the RPC-style
call will "feel" like you're executing the service work
locally, when in reality the client is actually making API
calls to the server running on port 9000.

```go
package main

import (
    "context"
    "fmt"

    "github.com/your/project/calc"
    calcgen "github.com/your/project/calc/gen"
)

func main() {
    ctx := context.Background()
    client := calcgen.CalculatorServiceClient("http://localhost:9000")

    add, err := client.Add(ctx, &calc.AddRequest{A:5, B:2})
    if err != nil {
        // handle error	
    }
    fmt.Println("5 + 2 =", add.Result)

    sub, err := client.Sub(ctx, &calc.SubRequest{A:5, B:2})
    if err != nil {
        // handle error	
    }
    fmt.Println("5 - 2 =", sub.Result)
}
```

Compile/run this program, and you should see the following output:

```
5 + 2 = 7
5 - 2 = 3
```
That's it!

For more examples of how to write services that let Abide take
care of the RPC/API boilerplate, take a look in the [example/](https://github.com/monadicstack/abide/tree/main/example)
directory of this repo.

## Adding Event-Driven Methods

RPC-style communication works for lots of scenarios, but sometimes
you want loosely-coupled workflows that fire when other operations
in the system complete. For instance, let's say that after a user
places an order, you want the system to send them an order confirmation
email as well as send them a coupon for a future order.

Abide makes it super easy to wire these events up. Here is what
your service interface would look like. And yes, you'd probably
put email-related operations in a different service, but we
just want to see how to wire up event-driven service invocation.
This setup works equally well if you split this up, but we'll
look at multiple service setups later.

```go
type OrderService interface {
    // PlaceOrder... places an order.
    //
    // HTTP 202
    // POST /orders
    PlaceOrder(context.Context, *PlaceOrderRequest) (*PlaceOrderResponse, error)

    // SendConfirmation sends an email confirmation to the user.
    //
    // ON OrderService.PlaceOrder
    SendConfirmation(context.Context, *SendConfirmationRequest) (*SendConfirmationResponse, error)

    // SendCoupon sends a 5% off similar item coupon to the user based on the order.
    //
    // ON OrderService.PlaceOrder
    SendCoupon(context.Context, *SendCouponRequest) (*SendCouponResponse, error)
}

type PlaceOrderResponse struct {
    OrderID  string
    UserID   string
    ItemIDs  string[]
    DateTime time.Time
}

type SendConfirmationRequest struct {
    OrderID string
    UserID  string
}

type SendCouponRequest struct {
    OrderID string
    UserID  string
}
```

The next section will go over Doc Options in more detail, but
just looking at the code, it should be somewhat obvious what
we're going for. We will invoke both "send" methods automatically be
any time there's a successful call to `PlaceOrder`.

When `PlaceOrder` finishes, Abide automatically publishes an
`OrderService.PlaceOrder` event with the response value.
Both `SendXXX` methods receive that event and build
their request structs automatically. They'll fill in `OrderID`
and `UserID`, but they'll just ignore `ItemIDs` and `DateTime`
because they don't have equivalent fields for those.

If that makes sense, notice that the only thing you did
differently than before was adding that line in the comments.
That's all the info that Abide needs to wire that behavior up for you!

There is one more one-line change we need to make in order for
this to work. That's in `main()` when we set up our server.
Before we only told the server to listen for requests via an API
Gateway. Now we need to tell it to also listen for requests via
an Event Gateway:

```go
func main() {
    // Create the handler and service, exact same as before...
    orderHandler := orders.OrderServiceHandler{}
    orderService := ordersgen.NewOrderService(orderHandler)

    // Now, the service can accept requests via the HTTP API
    // OR from events wired up using the 'ON' doc option.
    server := runtime.New(orderService)
    server.Listen(apis.NewGateway(":9000"))
    server.Listen(events.NewGateway())
    server.Run()
}
```

Now you'd expect something like this when running this code:

```shell
curl -d '{...}' http://localhost:9000/orders
# curl response 
{
  "OrderID": "123",
  "UserID": "456",
  "ItemIDs": ["789"],
  "DateTime": "2022-12-17T17:00:23+00:00"
}
# and you should have 2 emails in your inbox
```

### Distributed Events Using NATS JetStream

The order example above works great if you're running everything
in one process as a monolith. By default, the Event Gateway uses an in-memory
event broker to publish and react to events fired by your services. 


If you want to write this as a
distributed system with multiple remote instances and services, however, you
will need some third party event broker to manage this. Abide ships with support
for using [NATS JetStream](https://docs.nats.io/nats-concepts/jetstream) out-of-the-box.

```go
import (
    // ... other imports ...
    "github.com/monadicstack/abide/eventsource/nats"
)

func main() {
    // Create the handler and service, exact same as before...
    orderHandler := orders.OrderServiceHandler{}
    orderService := ordersgen.NewOrderService(orderHandler)
	
    // Configure a NATS client to distribute events.
    natsBroker := nats.Broker(
        nats.WithAddress("nats://127.0.0.1:4222"),
        nats.WithMaxAge(24 * time.Hours),
    )
	
    // Tell the event gateway to use NATS instead of local queues.
    server := runtime.New(orderService)
    server.Listen(apis.NewGateway())
    server.Listen(events.NewGateway(
        events.WithBroker(natsBroker),	
    ))
    server.Run()
}
```

Now, you can run 20 different instances if you like, and the
events will be spread around to all of them rather than always being
handled by the instance that placed the order.

### A Word About "Consumer Groups"

If you were to run 20 instances of the `OrderService`, you're not going to
blast of 20 copies of each email. The NATS broker will create two Queue Groups
(consumer groups to use the Kafka-style terminology), one named "OrderService.SendConfirmation"
and another named "OrderService.SendCoupon". This means that when the place order event
fires, both groups will receive a copy of the event - BUT - only one of the 20 instances
will handle the event for the confirmation group and only one of the 20 instances
will handle the event for the coupon group. As a result, you can have as many loosely
coupled units of work fire while still scaling out your infrastructure.

## Doc Options: Custom URLs, Status, etc

Abide gives you a service/API that "just works" out of the
box. By default, endpoints follow a similar RPC/POST style used by lots of
other service libraries/frameworks.

You can, however customize the API routes for individual operations,
set a prefix for all routes in a service, and more using "Doc Options"...
worst Spider-Man villain ever.

Here's an example with most of the available options. They are all
independent, so you can specify a custom status without specifying
a custom route and so on.

```go
// CalculatorService provides some basic arithmetic operations.
//
// VERSION 0.1.3
// PATH /v1
type CalculatorService interface {
    // Add calculates the sum of A + B.
    //
    // HTTP 202
    // GET /sum/:A/:B
    Add(context.Context, *AddRequest) (*AddResponse, error)

    // Sub calculates the difference of A - B.
    //
    // GET /difference/:A/:B
    Sub(context.Context, *SubRequest) (*SubResponse, error)
	
    // CountCalls is NOT exposed in your HTTP API. It is, however
    // called after every single successful call to either Add or Sub
    // in this service. It will even fire when FixWord is called in
    // a completely different service!
    //
    // HTTP OMIT
    // ON CalculatorService.Add
    // ON CalculatorService.Sub
    // ON SpellingService.FixWord
    CountCalls(context.Context, *CountCallsRequest) (*CountCallsResponse, error)
}
```

#### Service: PATH

This prepends your custom value on every route in the API. It applies
to the standard `ServiceName.FunctionName` routes as well as custom routes
as we'll cover in a moment.

Your generated API and RPC clients will be auto-wired to use the prefix "v1" under the
hood, so you don't need to change your code any further. If you want
to hit the raw HTTP endpoints, however, here's how they look now:

```shell
curl -d '{"A":5, "B":2}' http://localhost:9000/v1/CalculatorService.Add
# {"Result":7}

curl -d '{"A":5, "B":2}' http://localhost:9000/v1/CalculatorService.Sub
# {"Result":3}
```

#### Function: GET/POST/PUT/PATCH/DELETE

You can replace the default `POST ServiceName.FunctionName` route for any
operation with the route of your choice. In the example, the path parameters `:A` and `:B`
will be bound to the equivalent A and B attributes on the request struct.

Here are the updated curl calls after we generate the new
gateway code. Notice it's also taking into account the service's PATH
prefix as well:

```shell
curl http://localhost:9000/v1/sum/5/2
# {"Result":7}
curl http://localhost:9000/v1/difference/5/2
# {"Result":3}
```

Use these options to your heart's content if you want your API
to feel more REST-ful instead of RPC-ful.

#### Function: HTTP {StatusCode}

This lets you have the API return a non-200 status code on success.
For instance, the Add function's route will return a `202 Accepted`
status when it responds with the answer instead of `200 OK`.

Since we didn't specify anything special for the Sub method, it
will continue to respond with `200 OK`, same as before.

#### Function: HTTP OMIT

Sometimes you want your service to be able to perform operations
that you don't want to expose to the outside world. Perhaps this
only fires asynchronously when some event fires (next section)
or it's just some private code you want to manually execute but not
allow external access.

If the operation has `HTTP OMIT`, Abide will not create an API
route for it. It will not appear in your OpenAPI docs or external
language clients (like JS and Dart). The method *will* still
appear in your Go client because we need to satisfy the service
interface, but you'll receive a 404 error if you attempt to
invoke it.

#### Function: ON {ServiceName.MethodName}

This is what we used in the previous section to allow services
to trigger workflow events. The format is always `ON ServiceName.MethodName`.
This is true even if you provide a custom HTTP API route. The
event name is ALWAYS the same no matter what.

As you can see in the example above, you can have as many `ON`
triggers as you want on a single method, and they do not even
need to be from the same service!

## Error Handling

By default, if your service call returns a non-nil error, the
resulting RPC/HTTP request will have a 500 status code. You
can, however, customize that status code to correspond to the type
of failure (e.g. 404 when something was not found).

The easiest way to do this is to just use Abide's `fail`
package when you encounter a failure case:

```go
import (
    "github.com/monadicstack/abide/fail"
)

func (svc UserService) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
    if req.ID == "" {
        return nil, fail.BadRequest("id is required")
    }
    user, err := svc.Repo.GetByID(req.ID)
    if err != nil {
    	return nil, err
    }
    if user == nil {
        return nil, fail.NotFound("user not found: %s", req.ID)
    }
    return &GetResponse{User: user}, nil
}
```

In this case, the caller will receive an HTTP 400 if they
didn't provide an id, a 404 if there is no user with that
id, and a 500 if any other type of error occurs.

#### Customizing Errors

While the error categories in the `fail` package are
probably good enough for most people, you can build your own
custom status-bound errors by simply having them implement the
`StatusCode() int` function:

```go
type RateLimitError struct {
    Limit int
}

func (err RateLimitError) StatusCode() int {
    return 429
}

func (err RateLimitError) Error() string {
    return fmt.Sprintf("limit of %d/sec exceeded", err.Limit)
}
```

Now when you implement a handler or middleware, you can simply
return your custom error type and have your API respond w/
a 429 instead of a generic 500 error:

```go
func (svc UserService) CreateToken(ctx context.Context, req *CreateTokenRequest) (*CreateTokenResponse, error) {
    if (svc.exceededLimit(ctx, 5)) {
        return nil, RateLimitError{Limit: 5}
    }
    return &CreateTokenResponse{Token: "Hello"}, nil
}

// Sample call:
// curl -XPOST http://localhost:9000/UserService.CreateToken
// {
//    "StatusCode": 429,
//    "Message": "limit of 5/sec exceeded"
// }
```

### Errors In Async Event Handlers

Handling errors in RPC calls is fairly easy. The clients that
Abide generate return the error. Simple.

When using the `ON Service.Method` option to trigger calls
based on events, you don't really have control over that code, so
we need to do something a little different to handle errors
that might occur during those asynchronous flows.

You can give the Event Gateway a callback function that Abide will
invoke any time an error occurs processing event-based service operations.

```go
func main() {
    // ...
    server.Listen(events.NewGateway(
        events.WithBroker(natsBroker),
        events.WithErrorHandler(handleEventError),
    ))
    server.Run()
}

func handleEventError(err error) {
    // Don't panic...
}
```

## Creating a JavaScript Client

The `abide` tool can actually generate a JS client that you
can add to your frontend code (or React Native mobile code)
to hide the complexity of making API calls to your backend
service. Without any plugins or fuss, we can create a JS client of the same
CalculatorService from earlier...

```shell
abide client calc/calculator_service.go --language=js
```

This will create the file `calculator_service.gen.client.js`
which you can include with your frontend codebase. Using it
should look similar to the Go client we saw earlier:

```js
import {CalculatorService} from 'lib/calculator_service.gen.client';

// The service client is a class that exposes all of the
// operations as 'async' functions that resolve with the
// result of the service call.
const service = new CalculatorService('http://localhost:9000');
const add = await service.Add({A:5, B:2});
const sub = await service.Sub({A:5, B:2});

// Should print:
// Add(5, 2) = 7
// Sub(5, 2) = 3
console.info('Add(5, 2) = ' + add.Result)
console.info('Sub(5, 2) = ' + sub.Result)
```

Another subtle benefit of using the generated client is that your
service/method documentation follows you in the generated code.
It's included in the file as JSDoc comments so your
documentation should be available to your IDE even when writing
your frontend code.

#### Node Support

Abide uses the `fetch` function to make the actual HTTP requests,
so if you are using Node 18+, you shouldn't need to do anything
special as `fetch` is now in the global scope. If that's the
case, ignore the next paragraph and subsequent sample code.

If you're using an older version of node or just really prefer
to use the classic `node-fetch` package, you can supply the
fetch implementation to use when constructing your client:

```js
const fetch = require('node-fetch');

const service = new CalculatorService('http://localhost:9000', {fetch});
const add = await service.Add({A:5, B:2});
const sub = await service.Sub({A:5, B:2});
```

## Creating a Dart/Flutter Client

Just like the JS client, Abide can create a Dart client that you can embed
in your Flutter apps so mobile frontends can consume your service.

```shell
abide client calc/calculator_service.go --language=dart
  or
abide client calc/calculator_service.go --language=flutter
```

This will create the file `calculator_service.gen.client.dart`. Add it
to your Flutter codebase, and it behaves very similarly to the JS client.

> The `HttpClient` from the standard `dart:io` package is NOT supported
> in Flutter web applications. To support Flutter mobile as well as web,
> Abide clients uses the [http](https://pub.dev/packages/http) package to
> make requests to the backend API. You'll need to add that to your
> pubspec for the following code to work:

```dart
import 'lib/calculator_service.gen.client.dart';

var service = CalculatorServiceClient("http://localhost:9000");
var add = await service.Add(AddRequest(A:5, B:2));
var sub = await service.Sub(SubRequest(A:5, B:2));

// Should print:
// Add(5, 2) = 7
// Sub(5, 2) = 3
print('Add(5, 2) = ${add.Result}');
print('Sub(5, 2) = ${sub.Result}');
```

## Middleware

You'll find that you frequently have work that you want to execute
before/after every single service invocation regardless of whether
it came from the API or some event. Abide uses continuation passing
functions similar to what you see in the most popular Go HTTP middleware
libraries.

```go
func main() {
    // Every service call to the CalculatorService will write to
    // the log and track how long it took.
    calcHandler := calc.CalculatorServiceHandler{}
    calcService := calcgen.NewCalculatorService(calcHandler,
        LogRequest,
        CollectTiming,
    )

    // No changes here...
    server := runtime.New(calcService)
    server.Listen(apis.NewGateway())
    server.Listen(events.NewGateway())
    server.Run()
}

func LogRequest(ctx context.Context, req any, next services.HandlerFunc) (any, error) {
    route := metadata.Route(ctx)
    fmt.Printf("Invoking %s.%s\n", route.ServiceName, route.Name)
    res, err := next(ctx, req)
    fmt.Printf(" > failed %v\n", err != nil)
    return res, err
}

func CollectTiming(ctx context.Context, req any, next services.HandlerFunc) (any, error) {
    start := time.Now()
    res, err := next(ctx, req)
    elapsed := time.Now().Sub(start)

    route := metadata.Route(ctx)
    metricsCollector.StoreTiming(route.QualifiedName(), elapsed)
    return res, err
}
```

#### HTTP Middleware

Most of your middleware should be done at the service level like
we have seen above: authorization, logging, observability, etc. They're
all things that are important regardless of whether we're servicing
an API call or an event.

One of Abide's primary goals is to make it so that you never have
to think about HTTP or transport code, but there are still times
when there's no getting around it. If you want to consume your
service in a web application, you're going to need to set up
CORS and that has to be done at the HTTP level.

Luckily, you can provide HTTP-level middleware when calling
`apis.NewGateway()`. The middleware function it expects are
compatible with [Negroni](https://github.com/urfave/negroni), so
you have an entire ecosystem of off-the-shelf handlers to plug in.

```go
func main() {
    // We'll still log and capture metrics on every call.
    calcHandler := calc.CalculatorServiceHandler{}
    calcService := calcgen.NewCalculatorService(calcHandler,
        LogRequest,
        CollectTiming,
    )

    server := runtime.New(calcService)
    server.Listen(apis.NewGateway(
        apis.WithMiddleware(
            negroni.NewLogger().ServeHTTP,
            cors.New().ServeHTTP,
            gzip.New().ServeHTTP,
        ),
    ))
    server.Listen(events.NewGateway())
    server.Run()
}
```

## Metadata

When you make an RPC call from Service A to Service B, values
stored on the `context.Context` will NOT be available to you when
are in Service B's handler. There are
instances, however, where it's useful to have data follow
every hop from service to service; trace ids, authorization, etc.

Abide uses the `metadata` package to store all manner of values for
the entire request; even if that request hits multiple services.

### Metadata: Authorization

You probably want your services to have some level of access control,
so incoming HTTP calls will likely have the "Authorization" header set.
Your middleware and handler functions don't work at the HTTP level, so the
`metadata` package captures that and makes it available for you.

```go
func CheckAdminMiddleware(ctx context.Context, req any, next services.Handler) (any, error) {
    auth := metadata.Authorization(ctx)
    if auth != "Bearer 12345" {
        return nil, fail.PermissionDenied("admins only")
    }
    return next(ctx, req)
}
```

Ignore the horrifyingly bad security - ultimately the service
should behave like this:

```shell
curl -H "Authorization: Guest" http://localhost:9000/AdminService.DropTables
{"StatusCode":403, "Message": "admins only"}

curl -H "Authorization: Bearer 12345" http://localhost:9000/AdminService.DropTables
{"TablesDropped":42}
```

#### Supplying Authorization Credentials

In the previous example we assumed that there was some primordial HTTP request with
an Authorization header. Since authorization is just a value stored
on the context, you can supply them fairly easily when using
the auto-generated service clients - again, the goal is to avoid
worrying about the transport layer:

```go
// Go: Using the generated service client
client := admingen.AdminServiceClient("http://localhost:9000")
ctx = metadata.WithAuthorization(ctx, "Bearer 12345")
client.DropTables(ctx, &admin.DropTablesRequest{Name:"*"})
```

```js
// JS: Using the generated service client
const client = new AdminClient('...');
const req = {Name: '*'};
client.DropTables(req, {authorization: 'Bearer 12345'});
```

```dart
// Dart: Using the generated service client
var client = AdminClient('...');
var req = DropTablesRequest(Name: '*');
client.DropTables(req, authorization: 'Token 12345');
```

### Metadata: Trace ID

If you ever want to be able to debug/observe behaviors in your
system, you'll need a consistent request/trace id to tie back
to every operation. For instance, if you place an order and then
that triggers 4 other operations (emails, analytics, etc.), Abide
manages a trace id that will be the same across all of those
related operations.

> Abide will honor any X-Request-ID header it receives,
so if your service is behind a load balancer or some proxy that
generates that HTTP header, that is the Trace ID that Abide will use.
If not, Abide will generate a unique value for you so that you always
have a meaningful Trace ID.

```go
func (svc FooService) Foo(ctx context.Context, req *FooRequest) (*FooResponse, error) {
    traceID := metadata.TraceID(ctx)
    fmt.Printf(">> Foo: %s\n", traceID)

    // By using the same context, the trace id is passed along.
    barServiceClient.Bar(ctx, &BarRequest{})
    // ...
}

func (svc BarService) Bar(ctx context.Context, req *BarRequest) (*BarResponse, error) {
    traceID := metadata.TraceID(ctx)
    fmt.Printf(">> Bar: %s\n", traceID)
    // ...
}

// The interface for BazService.Bar had this:
// ON BarService.Bar
func (svc BazService) Bar(ctx context.Context, req *BarRequest) (*BarResponse, error) {
    traceID := metadata.TraceID(ctx)
    fmt.Printf(">> Baz: %s\n", traceID)
    // ...
}
```
Here's the output of our console when we make that initial
call to the Foo operation:

```shell
# When your service receives an explicit request id:
curl -H "X-Request-ID: Hello12345" -XPOST http://localhost:9000/FooService.Foo
# Console output
>> Foo: Hello12345
>> Bar: Hello12345
>> Baz: Hello12345

# There's no explicit id, so we'll just create one:
curl -XPOST http://localhost:9000/FooService.Foo
# Console output
>> Foo: dGhpcyBpcyBhIHJlYWxseSBs
>> Bar: dGhpcyBpcyBhIHJlYWxseSBs
>> Baz: dGhpcyBpcyBhIHJlYWxseSBs
```

It doesn't matter how many hops your request takes or whether
they were RPC calls or event-based calls. Your trace id follows you.

### Metadata: Values

Although Abide manages some very specific fields with very specific
purposes, the `metadata` package lets you store a general purpose
map of values that you deem as important. Just like
authorization or trace ids, these values will be accessible by
subsequent service calls for the same request.

```go
func (svc ServiceA) Foo(ctx context.Context, req *FooRequest) (*FooResponse, error) {
    // "Hello" will NOT follow you when you call Bar(),
    // but "DontPanic" will. Notice that the metadata
    // value does not need to be a string like in gRPC.
    ctx = context.WithValue(ctx, "Hello", "World")
    ctx = metadata.WithValue(ctx, "DontPanic", Answer{Value: 42})

    serviceB.Bar(ctx, &BarRequest{})
}

func (b ServiceB) Bar(ctx context.Context, req *BarRequest) (*BarResponse, error) {
    valueA, okA := ctx.Value("Hello").(string)

    valueB := Answer{}
    okB = metadata.Value(ctx, "DontPanic", &b)
    
    // valueA  == ""               okA == false
    // valueB == Answer{Value:42}  okB == true
}

// Pretend that your ServiceC interface had this option on Baz:
// ON ServiceA.Foo
func (c ServiceC) Baz(ctx context.Context, req *BazRequest) (*BazResponse, error) {
    valueA, okA := ctx.Value("Hello").(string)

    valueB := Answer{}
    okB = metadata.Value(ctx, "DontPanic", &b)

    // valueA  == ""               okA == false
    // valueB == Answer{Value:42}  okB == true
}
```

If you're wondering why `metadata.Value()` looks more like
`json.Unarmsahl()` than `context.Value()`, it has to
do with a limitation of reflection in Go. When the values
are sent over the network from Service A to Service B/C, we
lose all type information. We need the type info `&b` gives
us in order to properly restore the original value, so Abide
follows the idiom established by many
of the decoders in the standard library.

## FAQs

### Why a separate repo/project? Why not do a Frodo version 2?

I hate Go the major versioning scheme for Go modules. I never understood the
widespread dislike of that choice until I ran into it myself. It's silly
having to make sure that people `go install` the URL that ends in `/v2` instead
of the more natural root package.

While Abide solves the same core issue that Frodo did, I changed the API
significantly to make multiservice deployments a first class citizen and
enable event-driven flows. Try as I might, I every attempt to fit event
driven stuff into Frodo's runtime code felt hacky and wrong. I needed a
version 2 but Go did us all dirty with versioning.

### Why does Abide only support NATS for event-driven flows?

Well, I had to start somewhere. NATS is written in Go, it's stupid simple to
set up, and satisfies most use cases, so it seemed like the natural way to go.
I may add support for Redis and (maybe) Kafka in the future, but for now NATS is
the only officially supported event driven client.
