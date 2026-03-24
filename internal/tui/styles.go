package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorAccent = lipgloss.Color("86")  // cyan
	colorDim    = lipgloss.Color("240")
	colorDir    = lipgloss.Color("75") // bold blue
	colorBinary = lipgloss.Color("238")

	styleDir = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorDir)

	styleFile = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	styleFileUnselected = lipgloss.NewStyle().
				Foreground(colorDim)

	styleBinary = lipgloss.NewStyle().
			Foreground(colorBinary).
			Italic(true)

	styleAccent = lipgloss.NewStyle().
			Foreground(colorAccent)

	styleStatusBar = lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(lipgloss.Color("250")).
			Padding(0, 1)

	stylePanelTitle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	stylePanelBorder = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("236"))

	styleExtChipOn = lipgloss.NewStyle().
			Background(colorAccent).
			Foreground(lipgloss.Color("0")).
			Padding(0, 1)

	styleExtChipOff = lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(colorDim).
			Padding(0, 1)
)
