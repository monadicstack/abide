package main

import (
	"fmt"
	"time"

	"github.com/monadicstack/abide/example/basic/calc"
	gen "github.com/monadicstack/abide/example/basic/calc/gen"
	"github.com/monadicstack/abide/services"
	"github.com/monadicstack/abide/services/gateways/apis"
	"github.com/monadicstack/abide/services/gateways/events"
)

func main() {
	fmt.Println("Initializing server")

	calcHandler := calc.CalculatorServiceHandler{}
	calcServer := gen.CalculatorServiceServer(calcHandler)

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
