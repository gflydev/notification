.PHONY: clean mod tidy build vet test test-race lint

clean:
	go clean -cache
	go clean -modcache
	go clean -testcache
	go mod tidy

mod:
	go list -m --versions

tidy:
	go mod tidy

build:
	go build ./...

vet:
	go vet ./...

test:
	go test -count=1 ./...

test-race:
	go test -race -count=1 ./...

lint: vet test-race
