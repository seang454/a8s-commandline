# A8S CLI Build And Install

## Purpose

This document explains how the A8S CLI source code becomes runnable binaries,
how to install those binaries on Windows, macOS, and Linux, and what users do
after installation.

## Build Flow

The source code lives in the repository:

```text
cmd/a8s/main.go        CLI entry point
internal/cli/          Cobra command tree and runtime
internal/auth/         Keycloak login, refresh, and logout
internal/api/          Backend HTTP client
internal/config/       Context and config loading
internal/operation/    YAML/JSON operation input
```

The build command compiles that source into an executable:

```text
source code -> go build -> dist/a8s.exe or dist/a8s-<os>-<arch>
```

Users do not need the source code after a release binary is built. They only
need the binary, a configuration file, and a Keycloak login.

## Build For Local Development

Use this when testing on your current Windows machine:

```powershell
New-Item -ItemType Directory -Force -Path .\dist
$env:GOCACHE = Join-Path (Get-Location) ".gocache"
go test ./...
go build -o dist\a8s.exe .\cmd\a8s
.\dist\a8s.exe version
```

The output is:

```text
dist/a8s.exe
```

That binary is only for local development/testing.

## Build Release Binaries For All Operating Systems

Use `make build-all` when `make` is available:

```bash
make build-all
```

Expected outputs:

```text
dist/a8s-windows-amd64.exe
dist/a8s-windows-arm64.exe
dist/a8s-linux-amd64
dist/a8s-linux-arm64
dist/a8s-darwin-amd64
dist/a8s-darwin-arm64
```

`darwin` means macOS.

On Windows without `make`, use PowerShell:

```powershell
New-Item -ItemType Directory -Force -Path .\dist
$env:GOCACHE = Join-Path (Get-Location) ".gocache"

$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o dist\a8s-windows-amd64.exe .\cmd\a8s
$env:GOOS="windows"; $env:GOARCH="arm64"; go build -o dist\a8s-windows-arm64.exe .\cmd\a8s

$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o dist\a8s-linux-amd64 .\cmd\a8s
$env:GOOS="linux"; $env:GOARCH="arm64"; go build -o dist\a8s-linux-arm64 .\cmd\a8s

$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o dist\a8s-darwin-amd64 .\cmd\a8s
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o dist\a8s-darwin-arm64 .\cmd\a8s

Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

## Which Binary Should Users Download?

| User machine | Binary |
|---|---|
| Windows Intel/AMD 64-bit | `a8s-windows-amd64.exe` |
| Windows ARM 64-bit | `a8s-windows-arm64.exe` |
| Linux Intel/AMD 64-bit | `a8s-linux-amd64` |
| Linux ARM 64-bit | `a8s-linux-arm64` |
| macOS Intel | `a8s-darwin-amd64` |
| macOS Apple Silicon | `a8s-darwin-arm64` |

## Install On Windows

From the project repository during development:

```powershell
.\scripts\install.ps1 -BinaryPath .\dist\a8s.exe
```

For release testing, install the release binary:

```powershell
.\scripts\install.ps1 -BinaryPath .\dist\a8s-windows-amd64.exe
```

The default install location is:

```text
%LOCALAPPDATA%\a8s\bin\a8s.exe
```

Add that folder to the user `PATH`:

```powershell
$installDir = "$env:LOCALAPPDATA\a8s\bin"
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")

if ($userPath -notlike "*$installDir*") {
  [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
}
```

Close and reopen PowerShell, then test:

```powershell
a8s version
```

## Install On Linux

Install the Linux binary:

```bash
./scripts/install.sh dist/a8s-linux-amd64
```

The default install location is:

```text
~/.local/bin/a8s
```

Add it to `PATH` if needed:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

To make that permanent, add the export line to `~/.bashrc`, `~/.zshrc`, or the
shell profile used by the user.

Test:

```bash
a8s version
```

## Install On macOS

For Apple Silicon:

```bash
./scripts/install.sh dist/a8s-darwin-arm64
```

For Intel macOS:

```bash
./scripts/install.sh dist/a8s-darwin-amd64
```

The default install location is:

```text
~/.local/bin/a8s
```

Add it to `PATH` if needed:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Test:

```bash
a8s version
```

For production macOS distribution, signed and notarized binaries are
recommended before public release.

## Configure After Install

The CLI needs a context config telling it which backend and Keycloak realm to
use.

During local development, set `A8S_CONFIG` to the repository example:

```powershell
[Environment]::SetEnvironmentVariable(
  "A8S_CONFIG",
  "D:\CSTADPreUniversityTraining\ITP\spring\a8s-commandline\.a8s.yaml",
  "User"
)
```

For Linux/macOS development:

```bash
export A8S_CONFIG="$PWD/.a8s.yaml"
```

For production users, the CLI should use the user config path:

```text
Windows: %APPDATA%\a8s\config.yaml
Linux:   ~/.config/a8s/config.yaml
macOS:   ~/Library/Application Support/a8s/config.yaml
```

The production config should point to the deployed backend:

```yaml
apiVersion: cli.a8s.io/v1alpha1
kind: Config

currentContext: production

contexts:
  production:
    server: https://api.autonomous-istad.com
    namespace: ns-user-01
    targetCluster: primary
    auth:
      issuer: https://keycloak.autonomous-istad.com/realms/a8s
      clientId: a8s-cli
      credentialKey: context:production
```

## First User Flow After Install

After install and config:

```bash
a8s auth login
a8s auth status
a8s project list
a8s workspace quota pricing
```

`a8s auth login` opens Keycloak in the browser. After the user logs in, the CLI
stores credentials locally and reuses them for later commands.

## Uninstall

Windows:

```powershell
Remove-Item -Force "$env:LOCALAPPDATA\a8s\bin\a8s.exe"
```

Linux/macOS:

```bash
rm -f "$HOME/.local/bin/a8s"
```

Uninstalling the binary does not automatically remove user config or stored
credentials.

