package main

import (
	"context"
	"fmt"
	"time"

	"github.com/monadicstack/abide/example/basic/calc"
	gen "github.com/monadicstack/abide/example/basic/calc/gen"
	"github.com/monadicstack/abide/metadata"
	"github.com/monadicstack/abide/services"
	"github.com/monadicstack/abide/services/gateways/apis"
	"github.com/monadicstack/abide/services/gateways/events"
)

func main() {
	fmt.Println("Initializing server")

	middleware := services.MiddlewareFuncs{
		func(ctx context.Context, req any, next services.HandlerFunc) (any, error) {
			fmt.Println(">>>> BEFORE MY FUNCTION, ", metadata.TraceID(ctx), metadata.Authorization(ctx))
			res, err := next(ctx, req)
			fmt.Println(">>>> AFTER MY FUNCTION", metadata.TraceID(ctx))
			return res, err
		},
	}

	calcHandler := calc.CalculatorServiceHandler{}
	calcServer := gen.CalculatorServiceServer(calcHandler, middleware...)

	server := services.NewServer(
		services.Listen(apis.NewGateway(":8080")),
		services.Listen(events.NewGateway()),
		services.Register(calcServer),
	)

	fmt.Println("Server running on http://localhost:8080")
	go server.ShutdownOnInterrupt(2 * time.Second)
	if err := server.Run(); err != nil {
		panic(err)
	}
}
