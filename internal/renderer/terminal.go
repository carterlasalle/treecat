package renderer

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/carterlasalle/treecat/internal/highlight"
	"github.com/carterlasalle/treecat/internal/scanner"
	"github.com/carterlasalle/treecat/internal/selector"
)

func renderTerminal(w io.Writer, state *selector.State, opts Options) error {
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
		if err := writeFileContent(w, node, opts); err != nil {
			return err
		}
	}
	return nil
}

func renderTree(w io.Writer, node *scanner.FileNode, prefix string, sortMode selector.SortMode) error {
	if node.Depth == 0 {
		if _, err := fmt.Fprintf(w, "%s/\n", node.Name); err != nil {
			return err
		}
	}
	children := sortedTreeChildren(node.Children, sortMode)
	for i, child := range children {
		last := i == len(children)-1
		connector := "├── "
		childPrefix := prefix + "│   "
		if last {
			connector = "└── "
			childPrefix = prefix + "    "
		}
		if _, err := fmt.Fprintf(w, "%s%s%s\n", prefix, connector, child.Name); err != nil {
			return err
		}
		if child.IsDir {
			if err := renderTree(w, child, childPrefix, sortMode); err != nil {
				return err
			}
		}
	}
	return nil
}

func sortedTreeChildren(children []*scanner.FileNode, mode selector.SortMode) []*scanner.FileNode {
	cp := make([]*scanner.FileNode, len(children))
	copy(cp, children)
	switch mode {
	case selector.SortSize:
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].IsDir != cp[j].IsDir {
				return cp[i].IsDir
			}
			return cp[i].Size > cp[j].Size
		})
	case selector.SortLines:
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].IsDir != cp[j].IsDir {
				return cp[i].IsDir
			}
			return cp[i].Lines > cp[j].Lines
		})
	case selector.SortExt:
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].IsDir != cp[j].IsDir {
				return cp[i].IsDir
			}
			if cp[i].Ext != cp[j].Ext {
				return cp[i].Ext < cp[j].Ext
			}
			return cp[i].Name < cp[j].Name
		})
	default:
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].IsDir != cp[j].IsDir {
				return cp[i].IsDir
			}
			return cp[i].Name < cp[j].Name
		})
	}
	return cp
}

func writeFileContent(w io.Writer, node *scanner.FileNode, opts Options) error {
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
			if _, err := fmt.Fprintf(w, "[binary — %d bytes, use --hex to view]\n", node.Size); err != nil {
				return err
			}
		}
		return nil
	}
	data, err := os.ReadFile(node.Path)
	if err != nil {
		return err
	}
	src := string(data)
	if opts.NoSyntax || opts.NoColor {
		if _, err := fmt.Fprint(w, src); err != nil {
			return err
		}
		if !strings.HasSuffix(src, "\n") {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		return nil
	}
	highlighted, err := highlight.File(node.Name, src, highlight.Options{Color: true})
	if err != nil || highlighted == "" {
		_, writeErr := fmt.Fprint(w, src)
		return writeErr
	}
	_, writeErr := fmt.Fprint(w, highlighted)
	return writeErr
}
