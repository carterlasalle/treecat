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

const longDescription = `treecat scans a directory, prints a directory tree, and exports the contents
of the selected files in terminal, markdown, or plain-text form.

It is designed for packaging codebase context for LLM prompts, code review,
and documentation handoffs while keeping output predictable for scripts.`

const examples = `  treecat
  treecat . -o md -f context.md
  treecat . -e .go -s lines
  treecat . --relative --tree-only
  treecat . -i
  treecat completion zsh > ~/.zsh/completions/_treecat
  treecat man > treecat.1`

// NewRootCmd builds the cobra root command, writing output to w.
func NewRootCmd(w io.Writer) *cobra.Command {
	var (
		outputFmt     string
		outputFile    string
		extensions    []string
		noIgnore      bool
		hidden        bool
		maxDepth      int
		maxSizeStr    string
		hexBinary     bool
		sortMode      string
		noTree        bool
		treeOnly      bool
		noColor       bool
		noSyntax      bool
		interactive   bool
		relativePaths bool
	)

	cmd := &cobra.Command{
		Use:               "treecat [path]",
		Short:             "Recursive directory tree + syntax-highlighted file contents",
		Long:              longDescription,
		Example:           examples,
		Version:           "dev",
		SilenceUsage:      true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
		DisableAutoGenTag: true,
		Args:              cobra.MaximumNArgs(1),
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
			format, err := parseFormat(outputFmt)
			if err != nil {
				return err
			}
			sort, err := parseSortMode(sortMode)
			if err != nil {
				return err
			}
			resolvedNoColor := noColor || os.Getenv("NO_COLOR") != ""

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
					Format:        format,
					HexBinary:     hexBinary,
					GitignorePath: gitignorePath,
					NoIgnore:      noIgnore,
					ShowHidden:    hidden,
					Extensions:    exts,
					MaxSize:       maxSizeBytes,
					SortMode:      sort,
					NoColor:       resolvedNoColor,
					NoSyntax:      noSyntax,
					RootPath:      abs,
					RelativePaths: relativePaths,
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
				Format:        format,
				NoColor:       resolvedNoColor,
				NoSyntax:      noSyntax,
				NoTree:        noTree,
				NoContent:     treeOnly,
				HexBinary:     hexBinary,
				RootPath:      abs,
				RelativePaths: relativePaths,
			})
		},
	}
	cmd.SetOut(w)
	cmd.SetErr(os.Stderr)
	cmd.SetVersionTemplate("{{.Version}}\n")
	cmd.AddCommand(newCompletionCmd(w), newManCmd(w), newVersionCmd(w))

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
	f.BoolVar(&noColor, "no-color", false, "disable ANSI colors (also honored via NO_COLOR)")
	f.BoolVar(&noSyntax, "no-syntax", false, "disable syntax highlighting")
	f.BoolVarP(&interactive, "interactive", "i", false, "launch interactive TUI")
	f.BoolVar(&relativePaths, "relative", false, "render file headers relative to the scanned root")

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

func parseFormat(s string) (renderer.Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "terminal", "term", "":
		return renderer.FormatTerminal, nil
	case "md", "markdown":
		return renderer.FormatMarkdown, nil
	case "txt", "text":
		return renderer.FormatText, nil
	default:
		return renderer.FormatTerminal, fmt.Errorf("invalid --output %q (expected terminal, md, or txt)", s)
	}
}

func parseSortMode(s string) (selector.SortMode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "name", "":
		return selector.SortName, nil
	case "size":
		return selector.SortSize, nil
	case "lines":
		return selector.SortLines, nil
	case "ext":
		return selector.SortExt, nil
	default:
		return selector.SortName, fmt.Errorf("invalid --sort %q (expected name, size, lines, or ext)", s)
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
