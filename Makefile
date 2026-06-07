BINARY_NAME=a8s
VERSION=0.1.0
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

build:
	go build -ldflags="-X github.com/yourname/a8s/pkg/version.Version=$(VERSION) -X github.com/yourname/a8s/pkg/version.BuildDate=$(BUILD_DATE)" -o $(BINARY_NAME) ./cmd/a8s

build-all:
	GOOS=linux   GOARCH=amd64 go build -ldflags="-X github.com/yourname/a8s/pkg/version.Version=$(VERSION) -X github.com/yourname/a8s/pkg/version.BuildDate=$(BUILD_DATE)" -o dist/$(BINARY_NAME)-linux-amd64 ./cmd/a8s
	GOOS=linux   GOARCH=arm64 go build -ldflags="-X github.com/yourname/a8s/pkg/version.Version=$(VERSION) -X github.com/yourname/a8s/pkg/version.BuildDate=$(BUILD_DATE)" -o dist/$(BINARY_NAME)-linux-arm64 ./cmd/a8s
	GOOS=darwin  GOARCH=amd64 go build -ldflags="-X github.com/yourname/a8s/pkg/version.Version=$(VERSION) -X github.com/yourname/a8s/pkg/version.BuildDate=$(BUILD_DATE)" -o dist/$(BINARY_NAME)-darwin-amd64 ./cmd/a8s
	GOOS=darwin  GOARCH=arm64 go build -ldflags="-X github.com/yourname/a8s/pkg/version.Version=$(VERSION) -X github.com/yourname/a8s/pkg/version.BuildDate=$(BUILD_DATE)" -o dist/$(BINARY_NAME)-darwin-arm64 ./cmd/a8s
	GOOS=windows GOARCH=amd64 go build -ldflags="-X github.com/yourname/a8s/pkg/version.Version=$(VERSION) -X github.com/yourname/a8s/pkg/version.BuildDate=$(BUILD_DATE)" -o dist/$(BINARY_NAME)-windows-amd64.exe ./cmd/a8s
	GOOS=windows GOARCH=arm64 go build -ldflags="-X github.com/yourname/a8s/pkg/version.Version=$(VERSION) -X github.com/yourname/a8s/pkg/version.BuildDate=$(BUILD_DATE)" -o dist/$(BINARY_NAME)-windows-arm64.exe ./cmd/a8s

test:
	go test ./...

generate-routes:
	go run ./scripts/generate-route-registry
	gofmt -w internal/cli/catalog/generated_routes.go
	gofmt -w internal/cli/features/*/routes_gen.go

generate-docs:
	go run ./scripts/generate-command-docs

checksums: build-all
	cd dist && sha256sum * > checksums.txt

clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

run:
	go run ./cmd/a8s $(ARGS)

tidy:
	go mod tidy

.PHONY: build build-all test clean run tidy generate-routes generate-docs checksums
