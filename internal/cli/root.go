package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/carterlasalle/treecat/internal/filter"
	"github.com/carterlasalle/treecat/internal/renderer"
	"github.com/carterlasalle/treecat/internal/scanner"
	"github.com/carterlasalle/treecat/internal/selector"
	"github.com/carterlasalle/treecat/internal/tui"
)

// NewRootCmd builds the cobra root command, writing output to w.
func NewRootCmd(w io.Writer) *cobra.Command {
	var (
		outputFmt   string
		outputFile  string
		extensions  []string
		noIgnore    bool
		hidden      bool
		maxDepth    int
		maxSizeStr  string
		hexBinary   bool
		sortMode    string
		noTree      bool
		treeOnly    bool
		noColor     bool
		noSyntax    bool
		interactive bool
	)

	cmd := &cobra.Command{
		Use:   "treecat [path]",
		Short: "Recursive directory tree + syntax-highlighted file contents",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) > 0 {
				root = args[0]
			}
			abs, err := filepath.Abs(root)
			if err != nil {
				return err
			}

			gitignorePath := ""
			if !noIgnore {
				gip := filepath.Join(abs, ".gitignore")
				if _, statErr := os.Stat(gip); statErr == nil {
					gitignorePath = gip
				}
			}

			var maxSizeBytes int64
			if maxSizeStr != "" {
				maxSizeBytes, err = parseSize(maxSizeStr)
				if err != nil {
					return fmt.Errorf("invalid --max-size: %w", err)
				}
			}

			exts := normaliseExts(extensions)
			sort := parseSortMode(sortMode)

			if interactive {
				source, err := scanner.Scan(abs, scanner.Options{
					MaxDepth: maxDepth,
					Hidden:   true,
				})
				if err != nil {
					return fmt.Errorf("scan: %w", err)
				}
				return tui.Run(source, tui.Options{
					Output:        w,
					Format:        parseFormat(outputFmt),
					HexBinary:     hexBinary,
					GitignorePath: gitignorePath,
					NoIgnore:      noIgnore,
					ShowHidden:    hidden,
					Extensions:    exts,
					MaxSize:       maxSizeBytes,
					SortMode:      sort,
				})
			}

			tree, err := scanner.Scan(abs, scanner.Options{
				MaxDepth: maxDepth,
				Hidden:   hidden,
			})
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			filtered := filter.Apply(tree, filter.Options{
				Extensions:      exts,
				LimitExtensions: len(exts) > 0,
				GitignorePath:   gitignorePath,
				NoIgnore:        noIgnore,
				MaxSize:         maxSizeBytes,
				IncludeHidden:   hidden,
			})

			state := selector.New(filtered)
			state.Sort(sort)

			out := w
			if outputFile != "" {
				f, err := os.Create(outputFile)
				if err != nil {
					return err
				}
				defer f.Close()
				out = f
			}

			return renderer.Render(out, state, renderer.Options{
				Format:    parseFormat(outputFmt),
				NoColor:   noColor,
				NoSyntax:  noSyntax,
				NoTree:    noTree,
				NoContent: treeOnly,
				HexBinary: hexBinary,
			})
		},
	}

	f := cmd.Flags()
	f.StringVarP(&outputFmt, "output", "o", "terminal", "output format: terminal|md|txt")
	f.StringVarP(&outputFile, "file", "f", "", "write output to file")
	f.StringSliceVarP(&extensions, "ext", "e", nil, "include extensions, e.g. .go,.ts")
	f.BoolVar(&noIgnore, "no-ignore", false, "disable .gitignore")
	f.BoolVar(&hidden, "hidden", false, "include hidden files")
	f.IntVar(&maxDepth, "max-depth", 0, "max directory depth (0=unlimited)")
	f.StringVar(&maxSizeStr, "max-size", "", "skip files larger than size, e.g. 1MB")
	f.BoolVar(&hexBinary, "hex", false, "hex dump binary files")
	f.StringVarP(&sortMode, "sort", "s", "name", "sort: name|size|lines|ext")
	f.BoolVar(&noTree, "no-tree", false, "skip directory tree header")
	f.BoolVar(&treeOnly, "tree-only", false, "print tree only, skip file contents")
	f.BoolVar(&noColor, "no-color", false, "disable ANSI colors")
	f.BoolVar(&noSyntax, "no-syntax", false, "disable syntax highlighting")
	f.BoolVarP(&interactive, "interactive", "i", false, "launch interactive TUI")

	return cmd
}

func normaliseExts(exts []string) []string {
	out := make([]string, 0, len(exts))
	for _, e := range exts {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		out = append(out, strings.ToLower(e))
	}
	return out
}

func parseFormat(s string) renderer.Format {
	switch strings.ToLower(s) {
	case "md", "markdown":
		return renderer.FormatMarkdown
	case "txt", "text":
		return renderer.FormatText
	default:
		return renderer.FormatTerminal
	}
}

func parseSortMode(s string) selector.SortMode {
	switch s {
	case "size":
		return selector.SortSize
	case "lines":
		return selector.SortLines
	case "ext":
		return selector.SortExt
	default:
		return selector.SortName
	}
}

func parseSize(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	multipliers := map[string]int64{
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
	}
	for suffix, mult := range multipliers {
		if strings.HasSuffix(s, suffix) {
			var n int64
			if _, err := fmt.Sscanf(strings.TrimSuffix(s, suffix), "%d", &n); err != nil {
				return 0, err
			}
			return n * mult, nil
		}
	}
	var n int64
	_, err := fmt.Sscanf(s, "%d", &n)
	return n, err
}
