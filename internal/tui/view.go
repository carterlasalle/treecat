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

	treeW := m.width * 2 / 5
	previewW := m.width - treeW - 4

	treeContent := renderTreePanel(m, treeW)
	previewContent := renderPreviewPanel(m, previewW)

	treePanel := stylePanelBorder.Width(treeW).Height(m.height - 4).Render(treeContent)
	previewPanel := stylePanelBorder.Width(previewW).Height(m.height - 4).Render(previewContent)

	top := lipgloss.JoinHorizontal(lipgloss.Top, treePanel, previewPanel)
	extBar := renderExtBar(m)
	statusBar := renderStatusBar(m)

	return lipgloss.JoinVertical(lipgloss.Left, top, extBar, statusBar)
}

func renderTreePanel(m Model, width int) string {
	var lines []string
	for i, fn := range m.flatNodes {
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
	return strings.Join(lines, "\n")
}

func renderPreviewPanel(m Model, _ int) string {
	if m.cursor >= len(m.flatNodes) {
		return ""
	}
	node := m.flatNodes[m.cursor].node
	title := stylePanelTitle.Render("Preview: " + node.Name) + "\n\n"

	if node.IsDir {
		return title + styleAccent.Render("(directory)")
	}

	if node.IsBinary {
		if m.showHex {
			data, _ := os.ReadFile(node.Path)
			return title + highlight.HexDump(data)
		}
		return title + styleBinary.Render(fmt.Sprintf("[binary — %s]", humanize.Bytes(uint64(node.Size)))) +
			"\n" + styleAccent.Render("Press H to toggle hex dump")
	}

	data, err := os.ReadFile(node.Path)
	if err != nil {
		return title + err.Error()
	}
	lines := strings.Split(string(data), "\n")
	maxLines := m.height - 8
	if maxLines < 1 {
		maxLines = 10
	}
	if len(lines) > maxLines {
		lines = append(lines[:maxLines], fmt.Sprintf("… (%d more lines)", len(lines)-maxLines))
	}
	return title + strings.Join(lines, "\n")
}

func renderExtBar(m Model) string {
	var chips []string
	for _, ext := range m.extOrder {
		if m.extSelected[ext] {
			chips = append(chips, styleExtChipOn.Render(ext))
		} else {
			chips = append(chips, styleExtChipOff.Render(ext))
		}
	}
	bar := strings.Join(chips, " ")
	return styleStatusBar.Width(m.width).Render("Ext: " + bar)
}

func renderStatusBar(m Model) string {
	stats := m.state.Stats()
	sortName := m.sortNames[m.sortMode]
	info := fmt.Sprintf("%d files · %s · sort:%s",
		stats.FileCount,
		humanize.Bytes(uint64(stats.TotalSize)),
		sortName,
	)
	hints := "↑↓ move  spc toggle  s sort  H hex  ctrl+g generate  q quit"
	gap := m.width - len(info) - len(hints) - 2
	if gap < 1 {
		gap = 1
	}
	return styleStatusBar.Width(m.width).Render(info + strings.Repeat(" ", gap) + hints)
}
