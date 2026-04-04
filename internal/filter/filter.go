package filter

import (
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"

	"github.com/carterlasalle/treecat/internal/scanner"
)

// Options controls which files are included.
type Options struct {
	Extensions      []string // empty = all extensions unless LimitExtensions is true
	LimitExtensions bool     // when true, only Extensions are kept (including empty => none)
	GitignorePath   string   // path to .gitignore file; empty = skip
	NoIgnore        bool     // if true, ignore the gitignore file
	MaxSize         int64    // bytes; 0 = unlimited
	IncludeHidden   bool     // if true, keep dot-prefixed files/dirs
}

// Apply returns a new tree with nodes excluded based on opts.
// Directories are kept only if they have any passing children, except the root.
func Apply(node *scanner.FileNode, opts Options) *scanner.FileNode {
	if len(opts.Extensions) > 0 {
		opts.LimitExtensions = true
	}

	var gi *gitignore.GitIgnore
	var giRoot string
	if opts.GitignorePath != "" && !opts.NoIgnore {
		gi, _ = gitignore.CompileIgnoreFile(opts.GitignorePath)
		giRoot = filepath.Dir(opts.GitignorePath)
	}
	extSet := make(map[string]struct{}, len(opts.Extensions))
	for _, e := range opts.Extensions {
		extSet[e] = struct{}{}
	}
	return applyNode(node, opts, gi, giRoot, extSet, true)
}

func applyNode(node *scanner.FileNode, opts Options, gi *gitignore.GitIgnore, giRoot string, exts map[string]struct{}, isRoot bool) *scanner.FileNode {
	if !opts.IncludeHidden && !isRoot && strings.HasPrefix(node.Name, ".") {
		return nil
	}

	// go-gitignore expects paths relative to the .gitignore directory.
	if gi != nil && giRoot != "" {
		rel, err := filepath.Rel(giRoot, node.Path)
		if err == nil && gi.MatchesPath(rel) {
			return nil
		}
	}

	clone := *node
	clone.Children = nil

	if node.IsDir {
		for _, c := range node.Children {
			if filtered := applyNode(c, opts, gi, giRoot, exts, false); filtered != nil {
				clone.Children = append(clone.Children, filtered)
			}
		}
		if !isRoot && len(clone.Children) == 0 {
			return nil
		}
		return &clone
	}

	if opts.MaxSize > 0 && node.Size > opts.MaxSize {
		return nil
	}
	if opts.LimitExtensions {
		if _, ok := exts[node.Ext]; !ok {
			return nil
		}
	}
	return &clone
}
