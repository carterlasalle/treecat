package scanner_test

import (
	"path/filepath"
	"testing"

	"github.com/carterlasalle/treecat/internal/scanner"
)

func TestScan_BasicTree(t *testing.T) {
	root, err := scanner.Scan("../../testdata/fixture", scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if root.Name != "fixture" {
		t.Fatalf("root name = %q, want fixture", root.Name)
	}
	if !root.IsDir {
		t.Fatal("root should be dir")
	}
	if len(root.Children) == 0 {
		t.Fatal("expected children")
	}
}

func TestScan_FileMetadata(t *testing.T) {
	root, err := scanner.Scan("../../testdata/fixture", scanner.Options{Hidden: true})
	if err != nil {
		t.Fatal(err)
	}
	node := findNode(root, "main.go")
	if node == nil {
		t.Fatal("main.go not found")
	}
	if node.Size == 0 {
		t.Error("size should be > 0")
	}
	if node.Lines == 0 {
		t.Error("lines should be > 0")
	}
	if node.Ext != ".go" {
		t.Errorf("ext = %q, want .go", node.Ext)
	}
	if node.IsBinary {
		t.Error("main.go should not be binary")
	}
}

func TestScan_BinaryDetection(t *testing.T) {
	root, err := scanner.Scan("../../testdata/fixture", scanner.Options{Hidden: true})
	if err != nil {
		t.Fatal(err)
	}
	node := findNode(root, "logo.png")
	if node == nil {
		t.Fatal("logo.png not found")
	}
	if !node.IsBinary {
		t.Error("logo.png should be detected as binary")
	}
}

func TestScan_HiddenExcludedByDefault(t *testing.T) {
	root, err := scanner.Scan("../../testdata/fixture", scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	if findNode(root, ".hidden") != nil {
		t.Error(".hidden dir should be excluded by default")
	}
}

func TestScan_HiddenIncludedWhenRequested(t *testing.T) {
	root, err := scanner.Scan("../../testdata/fixture", scanner.Options{Hidden: true})
	if err != nil {
		t.Fatal(err)
	}
	if findNode(root, ".hidden") == nil {
		t.Error(".hidden dir should be included when Hidden=true")
	}
}

func TestScan_MaxDepth(t *testing.T) {
	root, err := scanner.Scan("../../testdata/fixture", scanner.Options{MaxDepth: 1})
	if err != nil {
		t.Fatal(err)
	}
	srcNode := findNode(root, "src")
	if srcNode == nil {
		t.Fatal("src dir should be present at depth 1")
	}
	if len(srcNode.Children) != 0 {
		t.Error("src children should not be scanned at MaxDepth=1")
	}
}

func findNode(node *scanner.FileNode, name string) *scanner.FileNode {
	if filepath.Base(node.Name) == name {
		return node
	}
	for _, c := range node.Children {
		if found := findNode(c, name); found != nil {
			return found
		}
	}
	return nil
}
