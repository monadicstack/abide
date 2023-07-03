# Since our suites run as a single Go test (testify just makes it look like different
# tests), this timeout should be long enough to handle our slowest suite.
TEST_TIMEOUT=60s

#
# Builds the actual abide CLI executable.
#
build:
	@ \
	go build -o out/abide main.go

install: build
	@ \
 	echo "Overwriting go-installed version..." && \
 	cp out/abide $$GOPATH/bin/abide

#
# Runs the all of the test suites for the entire Abide module.
#
test: test-unit test-integration

#
# Runs the self-contained unit tests that don't require code generation or anything like that to run.
#
test-unit:
	@ \
	go test -count=1 -timeout $(TEST_TIMEOUT) -tags unit ./...

#
# Generates the clients for all of our supported languages (Go, JS, Dart) and runs tests on them
# to make sure that they all behave as expected. So not only can we generate them, but can we actually
# fetch data from the sample service and get the expected results back?
#
test-integration: generate
	@ \
	go test -v -count=1 -timeout $(TEST_TIMEOUT) -tags integration ./...

generate: build
	@ \
	go generate ./internal/testext/... && \
	mv ./internal/testext/gen/*.client.js ./generate/testdata/js/ && \
	mv ./internal/testext/gen/*.client.dart ./generate/testdata/dart/
