# A8S CLI Release Process

## Purpose

This document defines versioning, build, test, signing, packaging, publication, rollback, and support requirements for production CLI releases.

## Versioning

Use semantic versioning:

```text
MAJOR.MINOR.PATCH
```

Embed:

- CLI version
- Git commit
- build date
- Go version
- supported API version range

Example:

```bash
a8s version --client
a8s version --server
```

## Supported Platforms

Initial release matrix:

```text
windows/amd64
windows/arm64
linux/amd64
linux/arm64
darwin/amd64
darwin/arm64
```

## Required Release Gates

- `gofmt`/format checks pass
- `go vet` and lint pass
- unit and contract tests pass
- race tests pass where supported
- OpenAPI drift check passes
- endpoint-to-command coverage remains complete
- security scans pass
- integration smoke tests pass
- documentation generation is clean
- backend security production gate is approved

## Build Reproducibility

Use a release tool such as GoReleaser. Builds should:

- use pinned Go/tool versions
- set deterministic linker metadata
- avoid embedding credentials or local paths
- produce checksums
- generate an SBOM
- be reproducible where practical

## Signing

- sign release artifacts and checksum files
- use keyless signing or protected release keys
- publish verification instructions
- sign Windows binaries where distribution requires it

## Packaging and Distribution

Recommended channels:

- GitHub Releases or internal artifact registry
- Homebrew tap
- Scoop bucket
- Chocolatey package
- downloadable archives for Linux
- optional container image for CI usage

Each archive should include:

```text
a8s binary
LICENSE
README
checksums/verification guidance
shell completions
```

## Release Steps

1. Confirm changelog and compatibility notes.
2. Update version metadata.
3. Generate command docs and OpenAPI snapshot.
4. Run all release gates.
5. Tag the commit.
6. Build and sign artifacts.
7. Publish artifacts, checksums, SBOM, and release notes.
8. Run installation and smoke tests from published artifacts.
9. Promote package-manager manifests.
10. Monitor errors and rollback if required.

## Release Notes

Include:

- new commands and flags
- changed behavior
- deprecated or removed commands
- backend compatibility requirements
- security fixes
- migration instructions
- known issues

## Deprecation Policy

- mark deprecated commands in help output
- provide replacement commands
- keep aliases for at least one minor release when practical
- remove incompatible commands only in a major release

## Rollback

- retain previous signed artifacts
- document package-manager downgrade commands
- never overwrite an existing release tag
- publish an advisory when a release is withdrawn

## Post-Release Verification

Test published binaries:

```bash
a8s version --client
a8s completion powershell
a8s context list
a8s doctor
```

Run authenticated read-only smoke tests against a controlled environment.

## Production Security Gate

Before the first production release, verify:

- cluster and Kubernetes endpoints enforce authentication and ownership
- `/api/internal/**` uses service authentication
- Git integration authorization is reviewed
- WebSocket authentication is protected consistently
- `/api/admin/documentation/**` requires admin authorization
- error responses do not leak internal exceptions

