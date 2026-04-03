package renderer_test

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"

	"github.com/carterlasalle/treecat/internal/renderer"
	"github.com/carterlasalle/treecat/internal/scanner"
	"github.com/carterlasalle/treecat/internal/selector"
)

func makeState() *selector.State {
	root, _ := scanner.Scan("../../testdata/fixture", scanner.Options{})
	return selector.New(root)
}

func TestRender_Terminal_ContainsTree(t *testing.T) {
	var buf bytes.Buffer
	err := renderer.Render(&buf, makeState(), renderer.Options{
		Format:   renderer.FormatTerminal,
		NoColor:  true,
		NoSyntax: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, "Directory Structure") {
		t.Error("output should contain 'Directory Structure'")
	}
	if !strings.Contains(out, "main.go") {
		t.Error("output should contain main.go")
	}
}

func TestRender_Terminal_ContainsFileContents(t *testing.T) {
	var buf bytes.Buffer
	renderer.Render(&buf, makeState(), renderer.Options{
		Format:   renderer.FormatTerminal,
		NoColor:  true,
		NoSyntax: true,
	})
	out := buf.String()
	if !strings.Contains(out, "File:") {
		t.Error("output should contain file sections")
	}
	if !strings.Contains(out, "package main") {
		t.Error("output should contain file contents")
	}
}

func TestRender_Markdown_HasCodeFences(t *testing.T) {
	var buf bytes.Buffer
	renderer.Render(&buf, makeState(), renderer.Options{Format: renderer.FormatMarkdown})
	out := buf.String()
	if !strings.Contains(out, "```") {
		t.Error("markdown output should have code fences")
	}
	if !strings.Contains(out, "```go") {
		t.Error("markdown output should have language-tagged go fences")
	}
}

func TestRender_Text_NoANSI(t *testing.T) {
	var buf bytes.Buffer
	renderer.Render(&buf, makeState(), renderer.Options{Format: renderer.FormatText})
	out := buf.String()
	if strings.Contains(out, "\x1b[") {
		t.Error("text output should not contain ANSI codes")
	}
}

func TestRender_BinarySkipped(t *testing.T) {
	root, _ := scanner.Scan("../../testdata/fixture", scanner.Options{})
	s := selector.New(root)
	var buf bytes.Buffer
	renderer.Render(&buf, s, renderer.Options{Format: renderer.FormatText, NoColor: true})
	out := buf.String()
	if !strings.Contains(out, "[binary") {
		t.Error("binary files should show a [binary ...] note")
	}
}

func TestRender_HexDump(t *testing.T) {
	root, _ := scanner.Scan("../../testdata/fixture", scanner.Options{})
	s := selector.New(root)
	var buf bytes.Buffer
	renderer.Render(&buf, s, renderer.Options{
		Format:    renderer.FormatText,
		HexBinary: true,
	})
	out := buf.String()
	if !strings.Contains(out, "00000000") {
		t.Error("hex dump should contain offset for binary files")
	}
}

func TestRender_RelativePaths(t *testing.T) {
	var buf bytes.Buffer
	rootPath, err := filepath.Abs("../../testdata/fixture")
	if err != nil {
		t.Fatal(err)
	}
	err = renderer.Render(&buf, makeState(), renderer.Options{
		Format:        renderer.FormatText,
		NoColor:       true,
		RootPath:      rootPath,
		RelativePaths: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "File: main.go") {
		t.Fatal("expected relative file headers")
	}
}
