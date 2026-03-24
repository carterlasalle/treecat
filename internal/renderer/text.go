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
	plainOpts := opts
	plainOpts.NoColor = true
	plainOpts.NoSyntax = true

	if !opts.NoTree {
		fmt.Fprintf(w, "Directory Structure:\n\n")
		renderTree(w, state.Root, "", plainOpts)
		fmt.Fprintln(w)
	}
	if opts.NoContent {
		return nil
	}
	for _, node := range state.Selected() {
		fmt.Fprintf(w, "\n---\nFile: %s\n---\n\n", node.Path)
		if node.IsBinary {
			if opts.HexBinary {
				data, _ := os.ReadFile(node.Path)
				fmt.Fprint(w, highlight.HexDump(data))
			} else {
				fmt.Fprintf(w, "[binary — %d bytes]\n", node.Size)
			}
			continue
		}
		data, err := os.ReadFile(node.Path)
		if err != nil {
			return err
		}
		src := string(data)
		fmt.Fprint(w, src)
		if !strings.HasSuffix(src, "\n") {
			fmt.Fprintln(w)
		}
	}
	return nil
}
