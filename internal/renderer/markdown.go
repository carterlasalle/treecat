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
		fmt.Fprintf(w, "## Directory Structure\n\n")
		fmt.Fprint(w, "```\n")
		renderTree(w, state.Root, "", state.SortMode())
		fmt.Fprintf(w, "```\n\n")
	}
	if opts.NoContent {
		return nil
	}
	for _, node := range state.Selected() {
		lang := langFromExt(node.Ext)
		fmt.Fprintf(w, "### File: `%s`\n\n", node.Path)
		if node.IsBinary {
			if opts.HexBinary {
				data, err := os.ReadFile(node.Path)
				if err != nil {
					return err
				}
				fmt.Fprintf(w, "```\n%s```\n\n", highlight.HexDump(data))
			} else {
				fmt.Fprintf(w, "> [binary — %d bytes]\n\n", node.Size)
			}
			continue
		}
		data, err := os.ReadFile(node.Path)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "```%s\n%s", lang, string(data))
		if !strings.HasSuffix(string(data), "\n") {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "```\n\n")
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
