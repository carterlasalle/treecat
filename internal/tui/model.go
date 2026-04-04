package tui

import (
	"io"
	"os"
	"sort"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/carterlasalle/treecat/internal/filter"
	"github.com/carterlasalle/treecat/internal/renderer"
	"github.com/carterlasalle/treecat/internal/scanner"
	"github.com/carterlasalle/treecat/internal/selector"
)

// saveTarget controls where generated output is written.
type saveTarget int

const (
	saveTerminal saveTarget = iota // print to terminal after quit
	saveFile                       // write to file
	saveBoth                       // print to terminal AND write to file
)

// Options passed from CLI to TUI.
type Options struct {
	Output        io.Writer
	Format        renderer.Format
	HexBinary     bool
	GitignorePath string
	NoIgnore      bool
	ShowHidden    bool
	Extensions    []string
	MaxSize       int64
	SortMode      selector.SortMode
	NoColor       bool
	NoSyntax      bool
	RootPath      string
	RelativePaths bool
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
	source *scanner.FileNode
	state  *selector.State
	opts   Options

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

	showHex          bool
	showHidden       bool
	respectGitignore bool
	showHelp         bool
	done             bool

	selectedByPath  map[string]bool
	collapsedByPath map[string]bool

	// save dialog
	savePending bool
	saveTarget  saveTarget
	fileInput   textinput.Model
}

// NewModel creates a Model (exported for tests).
func NewModel(source *scanner.FileNode, opts Options) Model {
	return newModel(source, opts)
}

func newModel(source *scanner.FileNode, opts Options) Model {
	exts := collectExtensions(source)
	extOrder := make([]string, 0, len(exts))
	extSel := make(map[string]bool, len(exts))
	selectedExts := make(map[string]struct{}, len(opts.Extensions))
	for _, ext := range opts.Extensions {
		selectedExts[ext] = struct{}{}
	}
	limitExts := len(opts.Extensions) > 0
	for e := range exts {
		extOrder = append(extOrder, e)
		if limitExts {
			_, extSel[e] = selectedExts[e]
		} else {
			extSel[e] = true
		}
	}
	sort.Strings(extOrder)

	fi := textinput.New()
	fi.Placeholder = defaultOutputName(opts.Format)
	fi.CharLimit = 256
	fi.Width = 40

	m := Model{
		source:           source,
		opts:             opts,
		extOrder:         extOrder,
		extSelected:      extSel,
		sortMode:         opts.SortMode,
		sortNames:        []string{"name", "size", "lines", "ext"},
		showHidden:       opts.ShowHidden,
		respectGitignore: opts.GitignorePath != "" && !opts.NoIgnore,
		selectedByPath:   map[string]bool{},
		collapsedByPath:  map[string]bool{},
		fileInput:        fi,
	}
	m.applyFiltersAndRebuild("")
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
	if m.state != nil && m.state.Root != nil {
		m.flattenNode(m.state.Root, 0)
	}
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

func collectExtensions(root *scanner.FileNode) map[string]int {
	out := map[string]int{}
	var walk func(*scanner.FileNode)
	walk = func(node *scanner.FileNode) {
		if node == nil {
			return
		}
		if !node.IsDir {
			out[node.Ext]++
		}
		for _, child := range node.Children {
			walk(child)
		}
	}
	walk(root)
	return out
}

func (m *Model) enabledExtensions() []string {
	if !m.hasActiveExtensionFilter() {
		return nil
	}
	out := make([]string, 0, len(m.extOrder))
	for _, ext := range m.extOrder {
		if m.extSelected[ext] {
			out = append(out, ext)
		}
	}
	return out
}

func (m *Model) hasActiveExtensionFilter() bool {
	for _, ext := range m.extOrder {
		if !m.extSelected[ext] {
			return true
		}
	}
	return false
}

func (m *Model) currentCursorPath() string {
	if m.cursor >= 0 && m.cursor < len(m.flatNodes) {
		return m.flatNodes[m.cursor].node.Path
	}
	return ""
}

func (m *Model) syncStateMaps() {
	if m.state == nil || m.state.Root == nil {
		return
	}
	var walk func(*scanner.FileNode)
	walk = func(node *scanner.FileNode) {
		m.selectedByPath[node.Path] = node.Selected
		if node.IsDir {
			m.collapsedByPath[node.Path] = node.Collapsed
		}
		for _, child := range node.Children {
			walk(child)
		}
	}
	walk(m.state.Root)
}

func (m *Model) restoreNodeState(node *scanner.FileNode) {
	if selected, ok := m.selectedByPath[node.Path]; ok {
		node.Selected = selected
	}
	if node.IsDir {
		if collapsed, ok := m.collapsedByPath[node.Path]; ok {
			node.Collapsed = collapsed
		}
	}
	for _, child := range node.Children {
		m.restoreNodeState(child)
	}
}

func (m *Model) restoreCursor(path string) {
	if path != "" {
		for i, fn := range m.flatNodes {
			if fn.node.Path == path {
				m.cursor = i
				m.clampScroll()
				return
			}
		}
	}
	if len(m.flatNodes) == 0 {
		m.cursor = 0
		m.treeScroll = 0
		return
	}
	if m.cursor >= len(m.flatNodes) {
		m.cursor = len(m.flatNodes) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m.clampScroll()
}

func (m *Model) applyFiltersAndRebuild(cursorPath string) {
	m.syncStateMaps()
	if cursorPath == "" {
		cursorPath = m.currentCursorPath()
	}

	filtered := filter.Apply(m.source, filter.Options{
		Extensions:      m.enabledExtensions(),
		LimitExtensions: m.hasActiveExtensionFilter(),
		GitignorePath:   m.opts.GitignorePath,
		NoIgnore:        !m.respectGitignore,
		MaxSize:         m.opts.MaxSize,
		IncludeHidden:   m.showHidden,
	})
	state := selector.New(filtered)
	m.restoreNodeState(state.Root)
	state.Sort(m.sortMode)
	m.state = state
	m.rebuildFlat()
	m.restoreCursor(cursorPath)
}

func (m *Model) setAllExtensions(enabled bool) {
	for _, ext := range m.extOrder {
		m.extSelected[ext] = enabled
	}
}

func defaultOutputName(format renderer.Format) string {
	switch format {
	case renderer.FormatMarkdown:
		return "treecat.md"
	case renderer.FormatText:
		return "treecat.txt"
	default:
		return "treecat.txt"
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

// Run starts the Bubble Tea program and renders output when the user confirms.
func Run(source *scanner.FileNode, opts Options) error {
	m := newModel(source, opts)
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}
	fm, ok := finalModel.(Model)
	if !ok || !fm.done {
		return nil
	}

	renderOpts := renderer.Options{
		Format:        opts.Format,
		HexBinary:     opts.HexBinary,
		NoColor:       opts.NoColor,
		NoSyntax:      opts.NoSyntax,
		RootPath:      opts.RootPath,
		RelativePaths: opts.RelativePaths,
	}

	filePath := fm.fileInput.Value()
	if filePath == "" {
		filePath = defaultOutputName(opts.Format)
	}

	switch fm.saveTarget {
	case saveTerminal:
		return renderer.Render(opts.Output, fm.state, renderOpts)

	case saveFile:
		return renderToFile(filePath, fm.state, renderOpts)

	case saveBoth:
		if err := renderer.Render(opts.Output, fm.state, renderOpts); err != nil {
			return err
		}
		return renderToFile(filePath, fm.state, renderOpts)
	}
	return nil
}

func renderToFile(path string, state *selector.State, opts renderer.Options) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return renderer.Render(f, state, opts)
}
