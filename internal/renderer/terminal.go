package renderer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/carterlasalle/treecat/internal/highlight"
	"github.com/carterlasalle/treecat/internal/scanner"
	"github.com/carterlasalle/treecat/internal/selector"
)

func renderTerminal(w io.Writer, state *selector.State, opts Options) error {
	if !opts.NoTree {
		fmt.Fprintf(w, "Directory Structure:\n\n")
		renderTree(w, state.Root, "", opts)
		fmt.Fprintln(w)
	}
	if opts.NoContent {
		return nil
	}
	for _, node := range state.Selected() {
		fmt.Fprintf(w, "\n---\nFile: %s\n---\n\n", node.Path)
		if err := writeFileContent(w, node, opts); err != nil {
			return err
		}
	}
	return nil
}

func renderTree(w io.Writer, node *scanner.FileNode, prefix string, _ Options) {
	if node.Depth == 0 {
		fmt.Fprintf(w, "%s/\n", node.Name)
	}
	for i, child := range node.Children {
		last := i == len(node.Children)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if last {
			connector = "└── "
			childPrefix = prefix + "    "
		}
		fmt.Fprintf(w, "%s%s%s\n", prefix, connector, child.Name)
		if child.IsDir {
			renderTree(w, child, childPrefix, Options{})
		}
	}
}

func writeFileContent(w io.Writer, node *scanner.FileNode, opts Options) error {
	if node.IsBinary {
		if opts.HexBinary {
			data, err := os.ReadFile(node.Path)
			if err != nil {
				return err
			}
			fmt.Fprint(w, highlight.HexDump(data))
		} else {
			fmt.Fprintf(w, "[binary — %d bytes, use --hex to view]\n", node.Size)
		}
		return nil
	}
	data, err := os.ReadFile(node.Path)
	if err != nil {
		return err
	}
	src := string(data)
	if opts.NoSyntax || opts.NoColor {
		fmt.Fprint(w, src)
		if !strings.HasSuffix(src, "\n") {
			fmt.Fprintln(w)
		}
		return nil
	}
	highlighted, err := highlight.File(node.Name, src, highlight.Options{Color: true})
	if err != nil || highlighted == "" {
		fmt.Fprint(w, src)
		return nil
	}
	fmt.Fprint(w, highlighted)
	return nil
}
