BINARY_NAME=a8s
VERSION=0.1.0
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

build:
	go build -ldflags="-X github.com/yourname/a8s/pkg/version.Version=$(VERSION) -X github.com/yourname/a8s/pkg/version.BuildDate=$(BUILD_DATE)" -o $(BINARY_NAME) .

build-all:
	GOOS=linux   GOARCH=amd64 go build -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME)-windows-amd64.exe .

test:
	go test ./...

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

run:
	go run . $(ARGS)

tidy:
	go mod tidy

.PHONY: build build-all test clean run tidy
