# A8S CLI

A Go + Cobra CLI for the A8S platform.

The endpoint mapping and production design are complete. The current Go implementation is an early legacy subset and is being migrated to the resource-first architecture described in [docs/architecture.md](docs/architecture.md).

## Documentation

- [Documentation index](docs/README.md)
- [Backend API to CLI catalog](docs/backend-api-cli-catalog.md)
- [Architecture](docs/architecture.md)
- [Authentication](docs/authentication.md)
- [Configuration](docs/configuration.md)
- [Error contract](docs/error-contract.md)
- [Workflow contracts](docs/workflows.md)
- [Testing strategy](docs/testing-strategy.md)
- [Command-reference design](docs/command-reference.md)
- [OpenAPI compatibility](docs/openapi-compatibility.md)
- [Release process](docs/release-process.md)

## Quick Start

### 1. Configure
Copy `.a8s.yaml` to your home directory and set your API URL and token:
```bash
cp .a8s.yaml ~/.a8s.yaml
```

### 2. Build
```bash
make build
```

### 3. Use the current legacy subset
```bash
# List users
./a8s list users
./a8s list users --all
./a8s list users --output json

# List projects
./a8s list projects

# Create a user
./a8s create user --name "John Doe" --email "john@example.com"
./a8s create user --name "Admin User" --email "admin@example.com" --admin

# Delete a user
./a8s delete user --id "user-123"

# Version
./a8s version
```

The planned production syntax is resource-first:

```bash
a8s project list
a8s cluster create --file cluster.yaml --wait
a8s admin user create
a8s workspace quota purchase --plan premium --wait
```

## Current Legacy Configuration Priority
1. CLI flags (`--api-url`, `--token`)
2. Environment variables (`A8S_API_URL`, `A8S_API_TOKEN`)
3. Config file (`~/.a8s.yaml`)
4. Defaults (`http://localhost:8080`)

The production configuration will use named contexts and secure credential storage. See [docs/configuration.md](docs/configuration.md).

## Build for All Platforms
```bash
make build-all
# outputs: dist/a8s-linux-amd64, dist/a8s-darwin-arm64, dist/a8s-windows-amd64.exe
```
