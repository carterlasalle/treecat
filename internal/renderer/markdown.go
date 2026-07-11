package renderer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/carterlasalle/treecat/internal/highlight"
	"github.com/carterlasalle/treecat/internal/selector"
)

func renderMarkdown(w io.Writer, state *selector.State, opts Options) error {
	if !opts.NoTree {
		if _, err := fmt.Fprint(w, "## Directory Structure\n\n```\n"); err != nil {
			return err
		}
		if err := renderTree(w, state.Root, "", state.SortMode()); err != nil {
			return err
		}
		if _, err := fmt.Fprint(w, "```\n\n"); err != nil {
			return err
		}
	}
	if opts.NoContent {
		return nil
	}
	for _, node := range state.Selected() {
		lang := langFromExt(node.Ext)
		if _, err := fmt.Fprintf(w, "### File: `%s`\n\n", displayPath(node.Path, opts)); err != nil {
			return err
		}
		if node.IsBinary {
			if opts.HexBinary {
				data, err := os.ReadFile(node.Path)
				if err != nil {
					return err
				}
				if _, err := fmt.Fprintf(w, "```\n%s```\n\n", highlight.HexDump(data)); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(w, "> [binary — %d bytes]\n\n", node.Size); err != nil {
					return err
				}
			}
			continue
		}
		data, err := os.ReadFile(node.Path)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(w, "```%s\n%s", lang, string(data)); err != nil {
			return err
		}
		if !strings.HasSuffix(string(data), "\n") {
			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprint(w, "```\n\n"); err != nil {
			return err
		}
	}
	return nil
}

func langFromExt(ext string) string {
	m := map[string]string{
		".go":   "go",
		".ts":   "typescript",
		".tsx":  "tsx",
		".js":   "javascript",
		".jsx":  "jsx",
		".py":   "python",
		".rs":   "rust",
		".md":   "markdown",
		".json": "json",
		".yaml": "yaml",
		".yml":  "yaml",
		".sh":   "bash",
		".html": "html",
		".css":  "css",
		".sql":  "sql",
		".toml": "toml",
		".rb":   "ruby",
		".java": "java",
		".c":    "c",
		".cpp":  "cpp",
		".h":    "c",
	}
	if l, ok := m[ext]; ok {
		return l
	}
	return ""
}
