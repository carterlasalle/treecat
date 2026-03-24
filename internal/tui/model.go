package tui

import (
	"io"
	"sort"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/carterlasalle/treecat/internal/renderer"
	"github.com/carterlasalle/treecat/internal/scanner"
	"github.com/carterlasalle/treecat/internal/selector"
)

// Options passed from CLI to TUI.
type Options struct {
	Output    io.Writer
	Format    renderer.Format
	HexBinary bool
}

type panel int

const (
	panelTree panel = iota
	panelPreview
)

// flatNode is a visible row in the tree panel.
type flatNode struct {
	node  *scanner.FileNode
	depth int
}

// Model is the Bubble Tea model for the treecat TUI.
type Model struct {
	state *selector.State
	opts  Options

	width, height int
	focused       panel

	flatNodes     []*flatNode
	cursor        int
	treeScroll    int
	previewScroll int

	extOrder    []string
	extSelected map[string]bool

	sortMode  selector.SortMode
	sortNames []string

	showHex bool
	done    bool
}

// NewModel creates a Model (exported for tests).
func NewModel(state *selector.State, opts Options) Model {
	return newModel(state, opts)
}

func newModel(state *selector.State, opts Options) Model {
	exts := state.Extensions()
	extOrder := make([]string, 0, len(exts))
	extSel := make(map[string]bool, len(exts))
	for e := range exts {
		extOrder = append(extOrder, e)
		extSel[e] = true
	}
	sort.Strings(extOrder)

	m := Model{
		state:       state,
		opts:        opts,
		extOrder:    extOrder,
		extSelected: extSel,
		sortNames:   []string{"name", "size", "lines", "ext"},
	}
	m.rebuildFlat()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

// treePanelH returns the number of content rows visible in the tree panel.
// Total height minus: top border(1) + bottom border(1) + ext bar(1) + status bar(1) = 4.
func (m *Model) treePanelH() int {
	h := m.height - 4 - 2 // panel height minus borders
	if h < 1 {
		h = 10
	}
	return h
}

// rebuildFlat re-flattens the visible tree into m.flatNodes and re-clamps scroll.
func (m *Model) rebuildFlat() {
	m.flatNodes = nil
	m.flattenNode(m.state.Root, 0)
	m.clampScroll()
}

func (m *Model) flattenNode(node *scanner.FileNode, depth int) {
	if depth > 0 || !node.IsDir {
		m.flatNodes = append(m.flatNodes, &flatNode{node: node, depth: depth})
	}
	if node.IsDir && !node.Collapsed {
		children := sortedChildren(node.Children, m.sortMode)
		for _, c := range children {
			m.flattenNode(c, depth+1)
		}
	}
}

// sortedChildren returns a sorted copy of children without mutating the slice.
func sortedChildren(children []*scanner.FileNode, mode selector.SortMode) []*scanner.FileNode {
	cp := make([]*scanner.FileNode, len(children))
	copy(cp, children)
	switch mode {
	case selector.SortSize:
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].IsDir != cp[j].IsDir {
				return cp[i].IsDir
			}
			return cp[i].Size > cp[j].Size
		})
	case selector.SortLines:
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].IsDir != cp[j].IsDir {
				return cp[i].IsDir
			}
			return cp[i].Lines > cp[j].Lines
		})
	case selector.SortExt:
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].IsDir != cp[j].IsDir {
				return cp[i].IsDir
			}
			if cp[i].Ext != cp[j].Ext {
				return cp[i].Ext < cp[j].Ext
			}
			return cp[i].Name < cp[j].Name
		})
	default: // SortName
		sort.Slice(cp, func(i, j int) bool {
			if cp[i].IsDir != cp[j].IsDir {
				return cp[i].IsDir
			}
			return cp[i].Name < cp[j].Name
		})
	}
	return cp
}

// clampScroll ensures treeScroll keeps the cursor visible.
func (m *Model) clampScroll() {
	panelH := m.treePanelH()
	if m.cursor < m.treeScroll {
		m.treeScroll = m.cursor
	}
	if m.cursor >= m.treeScroll+panelH {
		m.treeScroll = m.cursor - panelH + 1
	}
	maxScroll := len(m.flatNodes) - panelH
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.treeScroll > maxScroll {
		m.treeScroll = maxScroll
	}
	if m.treeScroll < 0 {
		m.treeScroll = 0
	}
}

// Update handles messages. Full implementation in update.go.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return updateModel(m, msg)
}

// View renders the UI. Full implementation in view.go.
func (m Model) View() string {
	return renderView(m)
}

// Run starts the Bubble Tea program and renders output if the user confirms.
func Run(state *selector.State, opts Options) error {
	m := newModel(state, opts)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	fm, ok := finalModel.(Model)
	if !ok || !fm.done {
		return nil
	}
	// User pressed ctrl+g — render selected files to output.
	return renderer.Render(opts.Output, state, renderer.Options{
		Format:    opts.Format,
		HexBinary: opts.HexBinary,
	})
}
