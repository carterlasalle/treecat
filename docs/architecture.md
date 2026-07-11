# Architecture

Treecat keeps filesystem discovery, selection state, and output formatting separate so the command-line and TUI paths share the same behavior.

## Data flow

1. `internal/scanner` walks the requested root without following symlinks, records metadata, and marks binary files.
2. `internal/filter` applies ignore, hidden-file, extension, size, and depth rules.
3. `internal/selector` owns selected/collapsed state and deterministic sorting.
4. `internal/renderer` writes terminal, Markdown, or plain-text output to an `io.Writer`.
5. `internal/cli` validates flags and connects the non-interactive path.
6. `internal/tui` presents the same tree through Bubble Tea and returns the final selection for rendering.

## Design decisions

- **No symlink traversal:** scanner nodes reject symbolic links, preventing recursive cycles and reads outside the requested tree.
- **Streaming output:** renderers target `io.Writer`, allowing terminal and file output without constructing the complete result in memory.
- **Bounded binary inspection:** the scanner checks a small prefix for NUL bytes while collecting file metadata.
- **Stable automation surface:** CLI flags and output defaults follow semantic versioning; internal packages can evolve without becoming public API.
- **Reproducible releases:** CI tests the Go module, multiple operating systems, GoReleaser packaging, vulnerability scanning, and release artifacts before publication.
