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

	flatNodes []*flatNode
	cursor    int

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

// rebuildFlat re-flattens the visible tree into m.flatNodes.
// Full implementation lives in update.go.
func (m *Model) rebuildFlat() {
	m.flatNodes = nil
	m.flattenNode(m.state.Root, 0)
}

func (m *Model) flattenNode(node *scanner.FileNode, depth int) {
	if depth > 0 || !node.IsDir {
		m.flatNodes = append(m.flatNodes, &flatNode{node: node, depth: depth})
	}
	if node.IsDir && !node.Collapsed {
		for _, c := range node.Children {
			m.flattenNode(c, depth+1)
		}
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

// Run starts the Bubble Tea program.
func Run(state *selector.State, opts Options) error {
	m := newModel(state, opts)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
