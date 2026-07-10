.PHONY: mod tidy build vet test test-race lint

# List available module versions
mod:
	go list -m --versions

# Sync go.mod / go.sum with the source
tidy:
	go mod tidy

# Compile all packages
build:
	go build ./...

# Report suspicious constructs
vet:
	go vet ./...

# Run the test suite
test:
	go test -count=1 ./...

# Run the test suite with the race detector
test-race:
	go test -race -count=1 ./...

# Vet + race tests
lint: vet test-race
