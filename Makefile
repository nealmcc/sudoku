.PHONY: build

test:
	go test -race -count=1 ./...

build:
	go build

cover:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func coverage.out

coverall:
	go test -tags slow -covermode=count -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -html coverage.out
