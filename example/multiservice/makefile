#
# This runs the calculator service in a local process exposing it on port 9000.
#
run-server: generate
	go run cmd/main.go

#
# This runs the sample consumer/client program that utilizes the service client to make
# requests to the calculator service. This will fail if you do not have 'run-server'
# already running.
#
run-client:
	go run main.go

#
# Generates all of our RPC gateway/client artifacts from the service definitions.
#
generate:
#	abide server calc/calculator_service.go && \
#	abide client calc/calculator_service.go
