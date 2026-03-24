package highlight_test

import (
	"strings"
	"testing"

	"github.com/carterlasalle/treecat/internal/highlight"
)

func TestHighlight_GoFile(t *testing.T) {
	src := `package main

import "fmt"

func main() { fmt.Println("hello") }
`
	out, err := highlight.File("main.go", src, highlight.Options{Color: true})
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
	if !strings.Contains(out, "main") {
		t.Error("output should contain source text")
	}
}

func TestHighlight_Plain(t *testing.T) {
	src := "hello world"
	out, err := highlight.File("readme.md", src, highlight.Options{Color: false})
	if err != nil {
		t.Fatal(err)
	}
	if out != src {
		t.Errorf("plain output = %q, want %q", out, src)
	}
}

func TestHighlight_HexDump(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
	out := highlight.HexDump(data)
	if !strings.Contains(out, "00000000") {
		t.Error("hex dump should contain offset")
	}
	if !strings.Contains(out, "ff") {
		t.Error("hex dump should contain ff")
	}
}
