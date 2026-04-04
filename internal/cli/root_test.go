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

func TestCLI_RelativePaths(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--relative", "--no-color", "--no-syntax"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "File: main.go") {
		t.Fatalf("expected relative file path in output, got %q", out)
	}
}

func TestCLI_CompletionSubcommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"completion", "bash"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "treecat") {
		t.Fatal("expected generated completion script")
	}
}

func TestCLI_ManSubcommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"man"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, ".TH \"TREECAT\"") {
		t.Fatal("expected roff man page output")
	}
}

func TestCLI_VersionSubcommand(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.Version = "1.2.3"
	cmd.SetArgs([]string{"version"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(buf.String()); got != "1.2.3" {
		t.Fatalf("version output = %q, want 1.2.3", got)
	}
}

func TestCLI_InvalidFormat(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--output", "xml"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "invalid --output") {
		t.Fatalf("expected invalid output error, got %v", err)
	}
}

func TestCLI_InvalidSortMode(t *testing.T) {
	var buf bytes.Buffer
	cmd := cli.NewRootCmd(&buf)
	cmd.SetArgs([]string{"../../testdata/fixture", "--sort", "weird"})
	if err := cmd.Execute(); err == nil || !strings.Contains(err.Error(), "invalid --sort") {
		t.Fatalf("expected invalid sort error, got %v", err)
	}
}
