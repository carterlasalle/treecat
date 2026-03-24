package selector_test

import (
	"testing"

	"github.com/carterlasalle/treecat/internal/scanner"
	"github.com/carterlasalle/treecat/internal/selector"
)

func makeTree() *scanner.FileNode {
	return &scanner.FileNode{
		Name: "root", IsDir: true,
		Children: []*scanner.FileNode{
			{Name: "a.go", Ext: ".go", Size: 100, Lines: 5},
			{Name: "b.go", Ext: ".go", Size: 200, Lines: 10},
			{Name: "c.md", Ext: ".md", Size: 50, Lines: 3},
			{Name: "assets", IsDir: true, Children: []*scanner.FileNode{
				{Name: "logo.png", Ext: ".png", Size: 5000, IsBinary: true},
			}},
		},
	}
}

func TestSelector_AllSelectedByDefault(t *testing.T) {
	s := selector.New(makeTree())
	selected := s.Selected()
	if len(selected) != 4 { // a.go, b.go, c.md, logo.png
		t.Fatalf("expected 4 selected, got %d", len(selected))
	}
}

func TestSelector_Toggle(t *testing.T) {
	tree := makeTree()
	s := selector.New(tree)
	s.Toggle(tree.Children[0]) // deselect a.go
	for _, n := range s.Selected() {
		if n.Name == "a.go" {
			t.Error("a.go should be deselected")
		}
	}
}

func TestSelector_ToggleExt(t *testing.T) {
	s := selector.New(makeTree())
	s.ToggleExt(".go", false)
	for _, n := range s.Selected() {
		if n.Ext == ".go" {
			t.Error("all .go files should be deselected")
		}
	}
}

func TestSelector_Extensions(t *testing.T) {
	s := selector.New(makeTree())
	exts := s.Extensions()
	if exts[".go"] != 2 {
		t.Errorf(".go count = %d, want 2", exts[".go"])
	}
}

func TestSelector_Stats(t *testing.T) {
	s := selector.New(makeTree())
	stats := s.Stats()
	if stats.FileCount != 4 {
		t.Errorf("count = %d, want 4", stats.FileCount)
	}
	if stats.TotalSize != 5350 {
		t.Errorf("size = %d, want 5350", stats.TotalSize)
	}
}

func TestSelector_SortBySize(t *testing.T) {
	s := selector.New(makeTree())
	s.Sort(selector.SortSize)
	sel := s.Selected()
	if sel[0].Name != "logo.png" {
		t.Errorf("first = %q, want logo.png", sel[0].Name)
	}
}
