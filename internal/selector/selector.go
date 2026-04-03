package selector

import (
	"sort"

	"github.com/carterlasalle/treecat/internal/scanner"
)

// SortMode controls file ordering.
type SortMode int

const (
	SortName  SortMode = iota
	SortSize           // descending
	SortLines          // descending
	SortExt
)

// Stats summarises the current selection.
type Stats struct {
	FileCount  int
	TotalSize  int64
	TotalLines int64
}

// State holds mutable selection state over a FileNode tree.
type State struct {
	Root     *scanner.FileNode
	sortMode SortMode
}

// New creates a State with all files selected.
func New(root *scanner.FileNode) *State {
	s := &State{Root: root}
	s.setAll(root, true)
	return s
}

func (s *State) setAll(node *scanner.FileNode, selected bool) {
	node.Selected = selected
	for _, c := range node.Children {
		s.setAll(c, selected)
	}
}

// Toggle flips the selection of a single node (recursively for directories).
func (s *State) Toggle(node *scanner.FileNode) {
	node.Selected = !node.Selected
	if node.IsDir {
		s.setAll(node, node.Selected)
	}
}

// ToggleExt sets selection for all files with the given extension.
func (s *State) ToggleExt(ext string, selected bool) {
	s.toggleExtNode(s.Root, ext, selected)
}

func (s *State) toggleExtNode(node *scanner.FileNode, ext string, selected bool) {
	if !node.IsDir && node.Ext == ext {
		node.Selected = selected
	}
	for _, c := range node.Children {
		s.toggleExtNode(c, ext, selected)
	}
}

// Selected returns all selected non-directory nodes in traversal order.
func (s *State) Selected() []*scanner.FileNode {
	var out []*scanner.FileNode
	s.collectSelected(s.Root, &out)
	s.applySortMode(&out)
	return out
}

func (s *State) collectSelected(node *scanner.FileNode, out *[]*scanner.FileNode) {
	if !node.IsDir && node.Selected {
		*out = append(*out, node)
	}
	for _, c := range node.Children {
		s.collectSelected(c, out)
	}
}

func (s *State) applySortMode(nodes *[]*scanner.FileNode) {
	switch s.sortMode {
	case SortSize:
		sort.Slice(*nodes, func(i, j int) bool {
			return (*nodes)[i].Size > (*nodes)[j].Size
		})
	case SortLines:
		sort.Slice(*nodes, func(i, j int) bool {
			return (*nodes)[i].Lines > (*nodes)[j].Lines
		})
	case SortExt:
		sort.Slice(*nodes, func(i, j int) bool {
			return (*nodes)[i].Ext < (*nodes)[j].Ext
		})
	}
}

// Sort sets the sort mode.
func (s *State) Sort(mode SortMode) { s.sortMode = mode }

// SortMode returns the current sort mode.
func (s *State) SortMode() SortMode { return s.sortMode }

// Extensions returns a map of extension → file count across the whole tree.
func (s *State) Extensions() map[string]int {
	out := map[string]int{}
	s.collectExts(s.Root, out)
	return out
}

func (s *State) collectExts(node *scanner.FileNode, out map[string]int) {
	if !node.IsDir {
		out[node.Ext]++
	}
	for _, c := range node.Children {
		s.collectExts(c, out)
	}
}

// Stats returns aggregate counts for currently selected files.
func (s *State) Stats() Stats {
	var st Stats
	for _, n := range s.Selected() {
		st.FileCount++
		st.TotalSize += n.Size
		st.TotalLines += n.Lines
	}
	return st
}
