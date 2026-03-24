package tui_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"

	"github.com/carterlasalle/treecat/internal/renderer"
	"github.com/carterlasalle/treecat/internal/scanner"
	"github.com/carterlasalle/treecat/internal/selector"
	"github.com/carterlasalle/treecat/internal/tui"
)

func makeModel(t *testing.T) tui.Model {
	t.Helper()
	root, err := scanner.Scan("../../testdata/fixture", scanner.Options{})
	if err != nil {
		t.Fatal(err)
	}
	state := selector.New(root)
	return tui.NewModel(state, tui.Options{Format: renderer.FormatTerminal})
}

func TestTUI_RendersWithoutPanic(t *testing.T) {
	m := makeModel(t)
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updated, _ := m.Update(msg)
	view := updated.(tui.Model).View()
	if view == "" {
		t.Error("view should not be empty after window size message")
	}
}

func TestTUI_QuitOnQ(t *testing.T) {
	m := makeModel(t)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 40))
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))
}
