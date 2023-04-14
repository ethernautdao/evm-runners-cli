package tui

import (
	"fmt"
	"strings"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethernautdao/evm-runners-cli/internal/config"
)

type levelListModel struct {
	Levels           map[string]config.Level
	solves			 map[string]string
	Keys             []string
	Cursor           int
	Done             bool
	descriptionShown bool
}

func (m *levelListModel) Init() tea.Cmd {
	m.Keys = make([]string, 0, len(m.Levels))
	for k := range m.Levels {
		m.Keys = append(m.Keys, k)
	}
	return nil
}

func (m *levelListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down":
			if m.Cursor < len(m.Levels)-1 {
				m.Cursor++
			}
		case "right":
			m.descriptionShown = true
		case "left":
			m.descriptionShown = false
		case "enter":
			m.Done = true
			return m, tea.Quit
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	return m, cmd
}

func (m *levelListModel) View() string {
	var sb strings.Builder

	header := "\n  #	  NAME	  	  SOLVES\n"
	headerSeparator := strings.Repeat("-", len(header)+10) + "\n"

	sb.WriteString(header)
	sb.WriteString(headerSeparator)

	for i, k := range m.Keys {
		l := m.Levels[k]
		if m.Cursor == i {
			sb.WriteString("> ")
		} else {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("%d	  %s	  	  %s\n", i+1, l.Contract, m.solves[l.Contract]))
		if m.Cursor == i && m.descriptionShown {
			sb.WriteString("\n" + "  " + l.Description + "\n\n")
		}
	}

	sb.WriteString("\n\x1b[90m↑ / ↓ - Navigate | ← / → - Toggle Description\x1b[0m")

	return sb.String()
}

func NewLevelList(Levels map[string]config.Level, solves map[string]string) *levelListModel {
	return &levelListModel{Levels: Levels, solves: solves}
}
