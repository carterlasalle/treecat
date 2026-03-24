package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/carterlasalle/treecat/internal/cli"
)

func TestCLI_DefaultRun(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--no-color", "--no-syntax"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Directory Structure") {
		t.Error("expected tree in output")
	}
	if !strings.Contains(out, "main.go") {
		t.Error("expected main.go in output")
	}
}

func TestCLI_ExtensionFilter(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--ext", ".md", "--no-color", "--no-syntax"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "package main") {
		t.Error("go file content should be excluded when filtering for .md")
	}
}

func TestCLI_MarkdownOutput(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--output", "md"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "```go") {
		t.Error("markdown output should contain go code fence")
	}
}
