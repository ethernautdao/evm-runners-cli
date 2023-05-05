package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type langListModel struct {
	Options []string
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
	s := "Do you want to use a template?\n\n"

	for i, option := range m.Options {
		// Add a ">" symbol before the selected option
		if i == m.Cursor {
			s += "> "
		} else {
			s += "  "
		}
		s += option + "\n"
	}

	s += "\n\x1b[90m↑ / ↓ - Navigate | ENTER - select language\x1b[0m"

	return s
}

func NewLangListModel() *langListModel {
	return &langListModel{
		Options: []string{"solidity", "vyper", "huff", "no template"},
		Cursor:  0,
	}
}
