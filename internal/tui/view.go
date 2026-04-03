package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dustin/go-humanize"

	"github.com/carterlasalle/treecat/internal/highlight"
)

func renderView(m Model) string {
	if m.width == 0 {
		return "Loading..."
	}

	mainHeight := m.height - 4
	if mainHeight < 3 {
		mainHeight = 3
	}

	var top string
	if m.width < 100 || m.height < 18 {
		panelWidth := maxInt(m.width, 20)
		panelContent := renderTreePanel(m, panelWidth-4)
		if m.focused == panelPreview {
			panelContent = renderPreviewPanel(m, panelWidth-4)
		}
		top = stylePanelBorder.Width(panelWidth).Height(mainHeight).Render(panelContent)
	} else {
		treeW := m.width * 2 / 5
		previewW := m.width - treeW - 4

		treeContent := renderTreePanel(m, treeW)
		previewContent := renderPreviewPanel(m, previewW)

		treePanel := stylePanelBorder.Width(treeW).Height(mainHeight).Render(treeContent)
		previewPanel := stylePanelBorder.Width(previewW).Height(mainHeight).Render(previewContent)
		top = lipgloss.JoinHorizontal(lipgloss.Top, treePanel, previewPanel)
	}

	var line1, line2 string
	if m.savePending {
		line1 = renderSaveBar(m)
		line2 = renderSaveFileLine(m)
	} else if m.showHelp {
		line1 = renderHelpBar(m)
		line2 = renderStatusBar(m)
	} else {
		line1 = renderExtBar(m)
		line2 = renderStatusBar(m)
	}

	return lipgloss.JoinVertical(lipgloss.Left, top, line1, line2)
}

func renderTreePanel(m Model, width int) string {
	panelH := m.treePanelH()
	end := m.treeScroll + panelH
	if end > len(m.flatNodes) {
		end = len(m.flatNodes)
	}

	var lines []string
	for i := m.treeScroll; i < end; i++ {
		fn := m.flatNodes[i]
		node := fn.node
		indent := strings.Repeat("  ", fn.depth)

		checkbox := "[✓]"
		if !node.Selected {
			checkbox = "[ ]"
		}

		meta := ""
		if !node.IsDir {
			meta = " " + humanize.Bytes(uint64(node.Size))
			if node.Lines > 0 {
				meta += fmt.Sprintf(" %dL", node.Lines)
			}
			if node.IsBinary {
				meta += " [bin]"
			}
		}

		var line string
		if node.IsDir {
			arrow := "▼"
			if node.Collapsed {
				arrow = "▶"
			}
			name := fmt.Sprintf("%s %s", arrow, node.Name)
			line = styleDir.Render(fmt.Sprintf("%s%s %s", indent, checkbox, name))
		} else if node.IsBinary {
			line = styleBinary.Render(fmt.Sprintf("%s%s %s%s", indent, checkbox, node.Name, meta))
		} else if node.Selected {
			line = styleFile.Render(fmt.Sprintf("%s%s %s%s", indent, checkbox, node.Name, meta))
		} else {
			line = styleFileUnselected.Render(fmt.Sprintf("%s%s %s%s", indent, checkbox, node.Name, meta))
		}

		cursor := "  "
		if i == m.cursor {
			cursor = styleAccent.Render("> ")
		}
		line = cursor + line

		// Truncate to width
		runes := []rune(line)
		if len(runes) > width {
			line = string(runes[:width-1]) + "…"
		}
		lines = append(lines, line)
	}

	// Scroll indicators
	if m.treeScroll > 0 {
		indicator := styleAccent.Render(fmt.Sprintf("  ↑ %d above", m.treeScroll))
		lines = append([]string{indicator}, lines...)
		if len(lines) > panelH {
			lines = lines[:panelH]
		}
	}
	below := len(m.flatNodes) - end
	if below > 0 {
		indicator := styleAccent.Render(fmt.Sprintf("  ↓ %d below", below))
		if len(lines) < panelH {
			lines = append(lines, indicator)
		} else {
			lines[len(lines)-1] = indicator
		}
	}

	return strings.Join(lines, "\n")
}

func renderPreviewPanel(m Model, _ int) string {
	if len(m.flatNodes) == 0 {
		return stylePanelTitle.Render("Preview") + "\n\n" + styleAccent.Render("No files match the current filters")
	}
	if m.cursor >= len(m.flatNodes) {
		return ""
	}
	node := m.flatNodes[m.cursor].node

	if node.IsDir {
		title := stylePanelTitle.Render("Directory: "+node.Name) + "\n\n"
		return title + styleAccent.Render(fmt.Sprintf("%d children", len(node.Children)))
	}

	title := stylePanelTitle.Render("Preview: "+node.Name) + "\n\n"

	if node.IsBinary {
		if m.showHex {
			data, err := os.ReadFile(node.Path)
			if err != nil {
				return title + err.Error()
			}
			allLines := strings.Split(highlight.HexDump(data), "\n")
			return title + scrolledLines(allLines, m.previewScroll, m.height-8)
		}
		return title + styleBinary.Render(fmt.Sprintf("[binary — %s]", humanize.Bytes(uint64(node.Size)))) +
			"\n" + styleAccent.Render("Press x to toggle hex dump")
	}

	data, err := os.ReadFile(node.Path)
	if err != nil {
		return title + err.Error()
	}
	allLines := strings.Split(string(data), "\n")
	maxVisible := m.height - 8
	if maxVisible < 1 {
		maxVisible = 10
	}
	scrollInfo := ""
	if len(allLines) > maxVisible {
		scrollInfo = styleAccent.Render(
			fmt.Sprintf("  line %d/%d — tab+↑↓ to scroll\n", m.previewScroll+1, len(allLines)),
		)
	}
	return title + scrollInfo + scrolledLines(allLines, m.previewScroll, maxVisible)
}

// scrolledLines returns a window of lines starting at offset, clamped to available lines.
func scrolledLines(lines []string, offset, maxVisible int) string {
	if offset >= len(lines) {
		offset = len(lines) - 1
		if offset < 0 {
			offset = 0
		}
	}
	end := offset + maxVisible
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[offset:end], "\n")
}

func renderExtBar(m Model) string {
	var chips []string
	active := 0
	for _, ext := range m.extOrder {
		if m.extSelected[ext] {
			active++
			chips = append(chips, styleExtChipOn.Render(ext))
		} else {
			chips = append(chips, styleExtChipOff.Render(ext))
		}
	}
	bar := strings.Join(chips, " ")
	prefix := fmt.Sprintf("Ext: %d/%d active", active, len(m.extOrder))
	if bar != "" {
		prefix += " · " + bar
	}
	return styleStatusBar.Width(m.width).Render(prefix)
}

func renderStatusBar(m Model) string {
	stats := m.state.Stats()
	sortName := m.sortNames[m.sortMode]
	gitState := "git:on"
	if !m.respectGitignore {
		gitState = "git:off"
	}
	hiddenState := "hidden:on"
	if !m.showHidden {
		hiddenState = "hidden:off"
	}
	extState := "ext:all"
	if m.hasActiveExtensionFilter() {
		extState = "ext:filtered"
	}
	info := fmt.Sprintf("%d files · %s · sort:%s · %s · %s · %s", stats.FileCount, humanize.Bytes(uint64(stats.TotalSize)), sortName, gitState, hiddenState, extState)

	hints := "↑↓ move  ←→/enter fold  spc select  s sort  e ext  ? help  ctrl+g save  q quit"
	if m.focused == panelPreview {
		hints = "↑↓ scroll preview  tab tree  x hex  ? help  q quit"
	}
	if m.width < 100 || m.height < 18 {
		if m.focused == panelTree {
			hints = "tab preview  spc select  s sort  ? help"
		} else {
			hints = "tab tree  ↑↓ scroll  x hex  ? help"
		}
	}

	return styleStatusBar.Width(m.width).Render(composeStatusLine(m.width, info, hints))
}

func renderHelpBar(m Model) string {
	summary := "Help: space toggle · a/A select dir/all · s sort · e/E extensions · g gitignore · H hidden · x hex · ctrl+g save"
	if m.focused == panelPreview {
		summary = "Help: tab tree · ↑/↓ scroll · pgup/pgdn jump · x hex · ctrl+g save"
	}
	return styleStatusBar.Width(m.width).Render(truncateText(summary, m.width-2))
}

// renderSaveBar renders the top row of the save dialog.
func renderSaveBar(m Model) string {
	labels := []string{"Terminal", "File", "Both"}
	var chips []string
	for i, label := range labels {
		if saveTarget(i) == m.saveTarget {
			chips = append(chips, styleExtChipOn.Render("● "+label))
		} else {
			chips = append(chips, styleExtChipOff.Render("  "+label))
		}
	}
	bar := "Save: " + strings.Join(chips, "  ")
	hint := "Tab=cycle · Enter=save · Esc=cancel"
	return styleStatusBar.Width(m.width).Render(composeStatusLine(m.width, bar, hint))
}

// renderSaveFileLine renders the second row of the save dialog.
func renderSaveFileLine(m Model) string {
	if m.saveTarget == saveTerminal {
		return styleStatusBar.Width(m.width).Render("  → output will be printed to the terminal")
	}
	label := "  File: "
	if m.saveTarget == saveBoth {
		label = "  File (+ terminal): "
	}
	return styleStatusBar.Width(m.width).Render(label + m.fileInput.View())
}

func composeStatusLine(width int, left, right string) string {
	contentWidth := width - 2
	if contentWidth < 10 {
		contentWidth = 10
	}
	left = truncateText(left, contentWidth)
	right = truncateText(right, contentWidth)
	if lipgloss.Width(left)+1+lipgloss.Width(right) > contentWidth {
		right = truncateText(right, contentWidth/2)
	}
	if lipgloss.Width(left)+1+lipgloss.Width(right) > contentWidth {
		left = truncateText(left, contentWidth-lipgloss.Width(right)-1)
	}
	gap := contentWidth - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}
	return left + strings.Repeat(" ", gap) + right
}

func truncateText(s string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= width {
		return s
	}
	if width == 1 {
		return "…"
	}
	return string(runes[:width-1]) + "…"
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
