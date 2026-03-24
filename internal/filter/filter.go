package filter

import (
	"path/filepath"

	gitignore "github.com/sabhiram/go-gitignore"

	"github.com/carterlasalle/treecat/internal/scanner"
)

// Options controls which files are included.
type Options struct {
	Extensions    []string // empty = all extensions
	GitignorePath string   // path to .gitignore file; empty = skip
	NoIgnore      bool     // if true, ignore the gitignore file
	MaxSize       int64    // bytes; 0 = unlimited
}

// Apply returns a new tree with nodes excluded based on opts.
// Directories are kept if they have any passing children.
func Apply(node *scanner.FileNode, opts Options) *scanner.FileNode {
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
	return applyNode(node, opts, gi, giRoot, extSet)
}

func applyNode(node *scanner.FileNode, opts Options, gi *gitignore.GitIgnore, giRoot string, exts map[string]struct{}) *scanner.FileNode {
	// go-gitignore expects paths relative to the .gitignore directory
	if gi != nil && giRoot != "" {
		rel, err := filepath.Rel(giRoot, node.Path)
		if err == nil && gi.MatchesPath(rel) {
			return nil
		}
	}

	if node.IsDir {
		clone := *node
		clone.Children = nil
		for _, c := range node.Children {
			if filtered := applyNode(c, opts, gi, giRoot, exts); filtered != nil {
				clone.Children = append(clone.Children, filtered)
			}
		}
		return &clone
	}

	// File checks
	if opts.MaxSize > 0 && node.Size > opts.MaxSize {
		return nil
	}
	if len(exts) > 0 {
		if _, ok := exts[node.Ext]; !ok {
			return nil
		}
	}
	return node
}
