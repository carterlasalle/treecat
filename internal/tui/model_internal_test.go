package tui

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/carterlasalle/treecat/internal/renderer"
	"github.com/carterlasalle/treecat/internal/scanner"
)

func loadFixtureTree(t *testing.T) *scanner.FileNode {
	t.Helper()
	root, err := scanner.Scan("../../testdata/fixture", scanner.Options{Hidden: true})
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func hasNodeName(nodes []*flatNode, name string) bool {
	for _, node := range nodes {
		if node.node.Name == name {
			return true
		}
	}
	return false
}

func hasNodeExt(nodes []*flatNode, ext string) bool {
	for _, node := range nodes {
		if node.node.Ext == ext {
			return true
		}
	}
	return false
}

func TestModel_ToggleHideFiltersHiddenNodes(t *testing.T) {
	m := newModel(loadFixtureTree(t), Options{Format: renderer.FormatTerminal})
	if hasNodeName(m.flatNodes, ".hidden") {
		t.Fatal("hidden directory should be filtered by default")
	}

	updated, _ := handleKey(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("H")})
	m = updated.(Model)
	if !hasNodeName(m.flatNodes, ".hidden") {
		t.Fatal("hidden directory should be visible after toggling hidden files")
	}
}

func TestModel_ToggleGitignoreShowsIgnoredFiles(t *testing.T) {
	tempDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte("*.log\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "keep.txt"), []byte("keep\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "debug.log"), []byte("debug\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	root, err := scanner.Scan(tempDir, scanner.Options{Hidden: true})
	if err != nil {
		t.Fatal(err)
	}
	m := newModel(root, Options{Format: renderer.FormatTerminal, GitignorePath: filepath.Join(tempDir, ".gitignore")})
	if hasNodeName(m.flatNodes, "debug.log") {
		t.Fatal("gitignored file should be hidden by default")
	}

	updated, _ := handleKey(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("g")})
	m = updated.(Model)
	if !hasNodeName(m.flatNodes, "debug.log") {
		t.Fatal("gitignored file should appear when gitignore filtering is toggled off")
	}
}

func TestModel_ToggleExtFiltersCurrentExtension(t *testing.T) {
	m := newModel(loadFixtureTree(t), Options{Format: renderer.FormatTerminal, ShowHidden: true})
	for i, node := range m.flatNodes {
		if node.node.Name == "main.go" {
			m.cursor = i
			break
		}
	}
	if !hasNodeExt(m.flatNodes, ".go") {
		t.Fatal("fixture should start with visible .go files")
	}

	updated, _ := handleKey(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	m = updated.(Model)
	if hasNodeExt(m.flatNodes, ".go") {
		t.Fatal(".go files should be filtered after toggling their extension")
	}

	updated, _ = handleKey(m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("E")})
	m = updated.(Model)
	if !hasNodeExt(m.flatNodes, ".go") {
		t.Fatal("resetting extension filters should restore .go files")
	}
}

func TestModel_SaveDialogUsesFormatDefaultName(t *testing.T) {
	m := newModel(loadFixtureTree(t), Options{Format: renderer.FormatMarkdown})
	updated, _ := handleKey(m, tea.KeyMsg{Type: tea.KeyCtrlG})
	m = updated.(Model)
	if got := m.fileInput.Value(); got != "treecat.md" {
		t.Fatalf("save dialog default = %q, want treecat.md", got)
	}
}
