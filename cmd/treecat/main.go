package main

import (
	"fmt"
	"os"

	"github.com/carterlasalle/treecat/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd := cli.NewRootCmd(os.Stdout)
	cmd.Version = fmt.Sprintf("%s (commit %s, built %s)", version, commit, date)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
