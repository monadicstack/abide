package main

import (
	"fmt"
	"time"

	"github.com/monadicstack/abide/example/multiservice/dismissal"
	dismissgateway "github.com/monadicstack/abide/example/multiservice/dismissal/gen"
	"github.com/monadicstack/abide/example/multiservice/greetings"
	greetergateway "github.com/monadicstack/abide/example/multiservice/greetings/gen"
	"github.com/monadicstack/abide/services"
	"github.com/monadicstack/abide/services/gateways/apis"
)

func main() {
	fmt.Println("Initializing server")
	greeterService := greetings.GreeterServiceHandler{}
	dismissService := dismissal.DismissServiceHandler{}

	// Both service APIs will run in a single HTTP server running on localhost:8080.
	server := services.NewServer(
		services.Listen(apis.NewGateway(":8080")),
		services.Register(
			greetergateway.NewGreeterService(greeterService),
			dismissgateway.NewDismissService(dismissService),
		),
	)

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("Quick examples:")
	fmt.Println("  curl -XPOST -d '{\"Name\":\"Dude\"}' http://localhost:8080/GreeterService.Greet")
	fmt.Println("  curl -XPOST -d '{\"Name\":\"Walter\"}' http://localhost:8080/DismissService.Dismiss")

	// Fire up the API and shut down gracefully when we receive a SIGINT or SIGTERM signal.
	go server.ShutdownOnInterrupt(10 * time.Second)
	if err := server.Run(); err != nil {
		panic(err)
	}

	fmt.Println("Bye bye...")
}
