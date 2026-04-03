package renderer

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/carterlasalle/treecat/internal/selector"
)

// Format selects the output format.
type Format int

const (
	FormatTerminal Format = iota
	FormatMarkdown
	FormatText
)

// Options controls renderer behaviour.
type Options struct {
	Format        Format
	NoColor       bool
	NoSyntax      bool
	NoTree        bool // skip tree header
	NoContent     bool // skip file contents (tree only)
	HexBinary     bool // show hex dump for binary files
	RootPath      string
	RelativePaths bool
}

// Render writes the tree + file contents to w.
func Render(w io.Writer, state *selector.State, opts Options) error {
	switch opts.Format {
	case FormatMarkdown:
		return renderMarkdown(w, state, opts)
	case FormatText:
		return renderText(w, state, opts)
	default:
		return renderTerminal(w, state, opts)
	}
}

func displayPath(path string, opts Options) string {
	if !opts.RelativePaths || opts.RootPath == "" {
		return path
	}
	rel, err := filepath.Rel(opts.RootPath, path)
	if err != nil {
		return path
	}
	if rel == "." {
		return "."
	}
	return filepath.ToSlash(strings.TrimPrefix(rel, "./"))
}
