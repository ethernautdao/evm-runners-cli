package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"sort"
	"strings"
)

type levelListModel struct {
	Levels           map[string]utils.Level
	solves           map[string]string
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

	// Sort the keys based on the ID field in the utils.Level struct
	sort.Slice(m.Keys, func(i, j int) bool {
		return m.Levels[m.Keys[i]].ID < m.Levels[m.Keys[j]].ID
	})

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

	header := fmt.Sprintf("\n  #\t%-14s%s\t%s\t\t%s\n", "NAME", "SOLVES", "#1 GAS", "1# SIZE")
	headerSeparator := "\x1b[90m" + strings.Repeat("-", len(header)+28) + "\n" + "\x1b[0m"

	sb.WriteString(header)
	sb.WriteString(headerSeparator)

	for i, k := range m.Keys {
		l := m.Levels[k]
		if m.Cursor == i {
			sb.WriteString("> ")
		} else {
			sb.WriteString("  ")
		}
		sb.WriteString(fmt.Sprintf("%s\t%-14s%s\t\t%s\t%s\n", l.ID, strings.ToLower(l.Name), m.solves[l.Name], "placeholder", "placeholder"))
		if m.Cursor == i && m.descriptionShown {
			sb.WriteString("\n" + "\x1b[32m" + l.Description + "\x1b[0m" + "\n")
		}
	}

	sb.WriteString("\n\x1b[90m↑ / ↓ - Navigate | ← / → - Toggle Description | q to exit\x1b[0m")

	return sb.String()
}

func NewLevelList(Levels map[string]utils.Level, solves map[string]string) *levelListModel {
	return &levelListModel{Levels: Levels, solves: solves}
}
