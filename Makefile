BINARY=fcom

.PHONY: all build install clean test lint

all: build

build:
	go build -o $(BINARY) ./main.go

install:
	go install ./...

clean:
	rm -f $(BINARY)
	go clean

test:
	go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...

lint:
	golangci-lint run 