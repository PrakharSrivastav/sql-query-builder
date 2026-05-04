.PHONY: test test-integration test-all

test:
	go test ./...

test-integration:
	go test -tags=integration -v -timeout=300s ./...

test-all: test test-integration
