package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newCompletionCmd(w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "completion [bash|zsh|fish|powershell]",
		Short:     "Generate shell completion scripts",
		Long:      "Generate a completion script for the requested shell and write it to stdout.",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell", "pwsh"},
		RunE: func(cmd *cobra.Command, args []string) error {
			root := cmd.Root()
			switch strings.ToLower(args[0]) {
			case "bash":
				return root.GenBashCompletionV2(w, true)
			case "zsh":
				return root.GenZshCompletion(w)
			case "fish":
				return root.GenFishCompletion(w, true)
			case "powershell", "pwsh":
				return root.GenPowerShellCompletionWithDesc(w)
			default:
				return fmt.Errorf("unsupported shell %q (expected bash, zsh, fish, or powershell)", args[0])
			}
		},
	}
	return cmd
}

func newManCmd(w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "man",
		Short: "Generate a man page",
		Long:  "Generate a roff man page for treecat and write it to stdout.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			header := &doc.GenManHeader{
				Title:   "TREECAT",
				Section: "1",
				Manual:  "Treecat Manual",
				Source:  "treecat",
			}
			return doc.GenMan(cmd.Root(), header, w)
		},
	}
	return cmd
}

func newVersionCmd(w io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print treecat version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := fmt.Fprintln(w, cmd.Root().Version)
			return err
		},
	}
	return cmd
}
