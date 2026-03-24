package integration_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/carterlasalle/treecat/internal/cli"
)

func TestIntegration_FullPipeline_Terminal(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--no-color", "--no-syntax"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()

	if !strings.Contains(out, "Directory Structure") {
		t.Error("missing directory structure header")
	}
	if !strings.Contains(out, "main.go") {
		t.Error("missing main.go in tree")
	}
	if !strings.Contains(out, "package main") {
		t.Error("missing file contents")
	}
	if !strings.Contains(out, "[binary") {
		t.Error("binary files should be noted")
	}
}

func TestIntegration_FullPipeline_Markdown(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--output", "md"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "```go") {
		t.Error("markdown output should have go code fence")
	}
	if !strings.Contains(out, "## Directory Structure") {
		t.Error("markdown output should have directory structure header")
	}
}

func TestIntegration_ExtensionFilter(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--ext", ".md", "--no-color", "--no-syntax"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "package main") {
		t.Error("Go file contents should be excluded when filtering for .md only")
	}
	if !strings.Contains(out, "Fixture") {
		t.Error("README.md content should be present")
	}
}

func TestIntegration_HexDump(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--hex", "--no-color", "--no-syntax"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "00000000") {
		t.Error("hex dump should show offset for binary files")
	}
}

func TestIntegration_TreeOnly(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--tree-only", "--no-color"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if strings.Contains(out, "package main") {
		t.Error("tree-only should not include file contents")
	}
	if !strings.Contains(out, "main.go") {
		t.Error("tree-only should still show the tree")
	}
}

func TestIntegration_SortBySize(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--sort", "size", "--no-color", "--no-syntax"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	// Output should contain files — basic sanity check
	if !strings.Contains(out, "File:") {
		t.Error("sorted output should still contain file sections")
	}
}
