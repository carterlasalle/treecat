package filter_test

import (
	"path/filepath"
	"testing"

	"github.com/carterlasalle/treecat/internal/filter"
	"github.com/carterlasalle/treecat/internal/scanner"
)

const fixture = "../../testdata/fixture"

func scanFixture(t *testing.T) *scanner.FileNode {
	t.Helper()
	root, err := scanner.Scan(fixture, scanner.Options{Hidden: true})
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func TestFilter_ExtensionInclude(t *testing.T) {
	root := scanFixture(t)
	filtered := filter.Apply(root, filter.Options{Extensions: []string{".go"}})
	assertOnlyExt(t, filtered, ".go")
}

func TestFilter_Gitignore(t *testing.T) {
	root := scanFixture(t)
	// Add a synthetic .log node that should be gitignored
	abs, _ := filepath.Abs(fixture)
	logNode := &scanner.FileNode{
		Name: "debug.log",
		Ext:  ".log",
		Path: filepath.Join(abs, "debug.log"),
	}
	root.Children = append(root.Children, logNode)

	filtered := filter.Apply(root, filter.Options{GitignorePath: filepath.Join(abs, ".gitignore")})
	if findNode(filtered, "debug.log") != nil {
		t.Error("debug.log should be excluded by .gitignore")
	}
}

func TestFilter_MaxSize(t *testing.T) {
	root := scanFixture(t)
	bigNode := &scanner.FileNode{Name: "big.bin", Size: 10 * 1024 * 1024, Ext: ".bin"}
	root.Children = append(root.Children, bigNode)
	filtered := filter.Apply(root, filter.Options{MaxSize: 1 * 1024 * 1024})
	if findNode(filtered, "big.bin") != nil {
		t.Error("big.bin should be excluded by MaxSize")
	}
}

func TestFilter_NoIgnore(t *testing.T) {
	root := scanFixture(t)
	abs, _ := filepath.Abs(fixture)
	logNode := &scanner.FileNode{
		Name: "debug.log",
		Ext:  ".log",
		Path: filepath.Join(abs, "debug.log"),
	}
	root.Children = append(root.Children, logNode)
	filtered := filter.Apply(root, filter.Options{
		GitignorePath: filepath.Join(abs, ".gitignore"),
		NoIgnore:      true,
	})
	if findNode(filtered, "debug.log") == nil {
		t.Error("debug.log should be included when NoIgnore=true")
	}
}

func TestFilter_HiddenFiles(t *testing.T) {
	root := scanFixture(t)
	filtered := filter.Apply(root, filter.Options{IncludeHidden: false})
	if findNode(filtered, ".hidden") != nil {
		t.Error("hidden directories should be filtered when IncludeHidden=false")
	}
}

func TestFilter_PrunesEmptyDirs(t *testing.T) {
	root := &scanner.FileNode{
		Name:  "root",
		Path:  "/tmp/root",
		IsDir: true,
		Children: []*scanner.FileNode{{
			Name:  "empty",
			Path:  "/tmp/root/empty",
			IsDir: true,
		}},
	}
	filtered := filter.Apply(root, filter.Options{})
	if findNode(filtered, "empty") != nil {
		t.Error("empty directories should be pruned after filtering")
	}
}

func assertOnlyExt(t *testing.T, node *scanner.FileNode, ext string) {
	t.Helper()
	if node == nil {
		return
	}
	if node.IsDir {
		for _, c := range node.Children {
			assertOnlyExt(t, c, ext)
		}
		return
	}
	if node.Ext != ext {
		t.Errorf("unexpected file %q with ext %q", node.Name, node.Ext)
	}
}

func findNode(node *scanner.FileNode, name string) *scanner.FileNode {
	if node == nil {
		return nil
	}
	if node.Name == name {
		return node
	}
	for _, c := range node.Children {
		if found := findNode(c, name); found != nil {
			return found
		}
	}
	return nil
}
