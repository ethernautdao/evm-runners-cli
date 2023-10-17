package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type langListModel struct {
	Options []string
	Lang    []string
	Cursor  int
	Done    bool
}

func (m *langListModel) Init() tea.Cmd {
	return nil
}

func (m *langListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.Cursor--
			if m.Cursor < 0 {
				m.Cursor = len(m.Options) - 1
			}
		case "down":
			m.Cursor++
			if m.Cursor >= len(m.Options) {
				m.Cursor = 0
			}
		case "enter":
			m.Done = true
			return m, tea.Quit
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *langListModel) View() string {
	var sb strings.Builder

	if m.Done {
		return ""
	} else {
		sb.WriteString("Do you want to use a template?\n\n")
		sb.WriteString("\x1b[90m┌" + strings.Repeat("─", 16) + "┐\n\x1b[0m") // Top border of the box
		for i, option := range m.Options {
			// Add a ">" symbol before the selected option
			if i == m.Cursor {
				sb.WriteString("\x1b[90m│\x1b[0m> ")
			} else {
				sb.WriteString("\x1b[90m│\x1b[0m  ")
			}
			sb.WriteString(fmt.Sprintf("%-14s\x1b[90m│\x1b[0m\n", option))
		}
		sb.WriteString("\x1b[90m└" + strings.Repeat("─", 16) + "┘\n\x1b[0m") // Bottom border of the box
		sb.WriteString("\n\x1b[90m↑/↓ - Navigate | q to exit | ↩ to select \x1b[0m")
		return sb.String()
	}
}

func NewLangListModel() *langListModel {
	return &langListModel{
		Options: []string{"solidity", "yul", "vyper", "huff", "no template"},
		Lang:    []string{"sol", "yul", "vy", "huff", "no template"},
		Cursor:  0,
	}
}
