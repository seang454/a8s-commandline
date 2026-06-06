# A8S CLI

A production-grade CLI for your platform API, built with Go + Cobra.

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

### 3. Use
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

## Configuration Priority
1. CLI flags (`--api-url`, `--token`)
2. Environment variables (`A8S_API_URL`, `A8S_API_TOKEN`)
3. Config file (`~/.a8s.yaml`)
4. Defaults (`http://localhost:8080`)

## Build for All Platforms
```bash
make build-all
# outputs: dist/a8s-linux-amd64, dist/a8s-darwin-arm64, dist/a8s-windows-amd64.exe
```
