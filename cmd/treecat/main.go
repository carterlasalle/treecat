package main

import (
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("treecat %s (%s, built %s)\n", version, commit, date)
		os.Exit(0)
	}
	fmt.Println("treecat: not yet implemented")
}
