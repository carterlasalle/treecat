# Contributing

## Setup

Install the Go version declared in `go.mod`, then download the verified module graph:

```bash
go mod download
go mod verify
```

## Required checks

```bash
go mod tidy -diff
gofmt -w .
go vet ./...
go test ./... -race
golangci-lint run
goreleaser check
```

Run `git diff --exit-code` after formatting and module checks. New behavior should include tests at the narrowest useful layer; end-to-end CLI behavior belongs under `internal/integration`.

Do not commit compiled binaries or generated `dist/` contents. Release artifacts are built by GoReleaser from version tags.

## Pull requests

Explain the user impact, compatibility implications, and validation commands. Keep output-format changes backward compatible unless the pull request explicitly targets a major release.
