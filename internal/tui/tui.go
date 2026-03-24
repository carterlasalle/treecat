// Package tui provides the interactive Bubble Tea TUI for treecat.
// This file is a compile stub; the full implementation is in model.go,
// update.go, view.go, keys.go, and styles.go (added in Task 8+).
package tui

import (
	"io"

	"github.com/carterlasalle/treecat/internal/renderer"
	"github.com/carterlasalle/treecat/internal/selector"
)

// Options passed from CLI to TUI.
type Options struct {
	Output    io.Writer
	Format    renderer.Format
	HexBinary bool
}

// Run starts the interactive TUI. Replaced by full implementation in Task 8.
func Run(_ *selector.State, _ Options) error { return nil }
