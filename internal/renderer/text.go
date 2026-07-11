package renderer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/carterlasalle/treecat/internal/highlight"
	"github.com/carterlasalle/treecat/internal/selector"
)

func renderText(w io.Writer, state *selector.State, opts Options) error {
	if !opts.NoTree {
		if _, err := fmt.Fprint(w, "Directory Structure:\n\n"); err != nil {
			return err
		}
		if err := renderTree(w, state.Root, "", state.SortMode()); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}
	if opts.NoContent {
		return nil
	}
	for _, node := range state.Selected() {
		if _, err := fmt.Fprintf(w, "\n---\nFile: %s\n---\n\n", displayPath(node.Path, opts)); err != nil {
			return err
		}
		if node.IsBinary {
			if opts.HexBinary {
				data, err := os.ReadFile(node.Path)
				if err != nil {
					return err
				}
				if _, err := fmt.Fprint(w, highlight.HexDump(data)); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(w, "[binary — %d bytes]\n", node.Size); err != nil {
					return err
				}
			}
			continue
		}
		data, err := os.ReadFile(node.Path)
		if err != nil {
			return err
		}
		src := string(data)
		if _, err := fmt.Fprint(w, src); err != nil {
			return err
		}
		if !strings.HasSuffix(src, "\n") {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
	}
	return nil
}
