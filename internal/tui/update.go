package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func updateModel(m Model, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.clampScroll()
		return m, nil
	case tea.KeyMsg:
		return handleKey(m, msg)
	}
	return m, nil
}

func handleKey(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Save dialog intercepts all keys while open.
	if m.savePending {
		return handleSaveDialog(m, msg)
	}

	panelH := m.treePanelH()

	switch {
	case key.Matches(msg, keys.Quit):
		return m, tea.Quit

	case key.Matches(msg, keys.Up):
		if m.focused == panelPreview {
			if m.previewScroll > 0 {
				m.previewScroll--
			}
		} else {
			if m.cursor > 0 {
				m.cursor--
				m.previewScroll = 0
			}
		}

	case key.Matches(msg, keys.Down):
		if m.focused == panelPreview {
			m.previewScroll++
		} else {
			if m.cursor < len(m.flatNodes)-1 {
				m.cursor++
				m.previewScroll = 0
			}
		}

	case key.Matches(msg, keys.PageUp):
		if m.focused == panelPreview {
			m.previewScroll -= panelH / 2
			if m.previewScroll < 0 {
				m.previewScroll = 0
			}
		} else {
			m.cursor -= panelH
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.previewScroll = 0
		}

	case key.Matches(msg, keys.PageDown):
		if m.focused == panelPreview {
			m.previewScroll += panelH / 2
		} else {
			m.cursor += panelH
			if m.cursor >= len(m.flatNodes) {
				m.cursor = len(m.flatNodes) - 1
			}
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.previewScroll = 0
		}

	case key.Matches(msg, keys.Toggle):
		if m.cursor < len(m.flatNodes) {
			m.state.Toggle(m.flatNodes[m.cursor].node)
		}

	case key.Matches(msg, keys.Expand):
		if m.cursor < len(m.flatNodes) {
			n := m.flatNodes[m.cursor].node
			if n.IsDir {
				n.Collapsed = !n.Collapsed
			}
		}

	case key.Matches(msg, keys.Left):
		if m.cursor < len(m.flatNodes) {
			n := m.flatNodes[m.cursor].node
			if n.IsDir && !n.Collapsed {
				n.Collapsed = true
			}
		}

	case key.Matches(msg, keys.Right):
		if m.cursor < len(m.flatNodes) {
			n := m.flatNodes[m.cursor].node
			if n.IsDir && n.Collapsed {
				n.Collapsed = false
			}
		}

	case key.Matches(msg, keys.SelectAll):
		if m.cursor < len(m.flatNodes) {
			n := m.flatNodes[m.cursor].node
			if n.IsDir && len(n.Children) > 0 {
				target := !n.Children[0].Selected
				for _, c := range n.Children {
					c.Selected = target
				}
			}
		}

	case key.Matches(msg, keys.SelectAllR):
		m.state.Toggle(m.state.Root)

	case key.Matches(msg, keys.Sort):
		m.sortMode = (m.sortMode + 1) % 4
		m.state.Sort(m.sortMode)

	case key.Matches(msg, keys.Tab):
		if m.focused == panelTree {
			m.focused = panelPreview
		} else {
			m.focused = panelTree
		}

	case key.Matches(msg, keys.ToggleHex):
		m.showHex = !m.showHex
		m.previewScroll = 0

	case key.Matches(msg, keys.ToggleGit):
		// gitignore toggling requires rescan; noted for future enhancement

	case key.Matches(msg, keys.Confirm):
		// Open save dialog instead of immediately quitting.
		m.savePending = true
		m.saveTarget = saveTerminal
		m.fileInput.SetValue("output.md")
		m.fileInput.CursorEnd()
		return m, nil
	}

	m.rebuildFlat()
	return m, nil
}

// handleSaveDialog handles keys while the save destination picker is open.
// Tab cycles Terminal → File → Both. Enter confirms. Esc cancels.
// When File or Both is active, all other keys are forwarded to the text input.
func handleSaveDialog(m Model, msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.savePending = false
		m.fileInput.Blur()
		return m, nil

	case "enter":
		m.done = true
		return m, tea.Quit

	case "tab":
		m.saveTarget = (m.saveTarget + 1) % 3
		if m.saveTarget == saveTerminal {
			m.fileInput.Blur()
			return m, nil
		}
		cmd := m.fileInput.Focus()
		return m, cmd

	default:
		// Forward typing/editing keys to the file path input.
		if m.saveTarget != saveTerminal {
			var cmd tea.Cmd
			m.fileInput, cmd = m.fileInput.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}
