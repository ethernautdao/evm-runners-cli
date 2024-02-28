package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
)

type levelListModel struct {
	Levels           map[string]utils.Level
	solves           map[string]string
	submissions      map[string]string
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

	tableWidth := 75

	if m.Done {
		return ""
	} else {
		sb.WriteString("\x1b[90m┌" + strings.Repeat("─", tableWidth) + "┐\n\x1b[0m") // Top border of the box

		header := fmt.Sprintf("\x1b[90m│\x1b[0m  #\t%-20s%-14s%-14s%-20s\x1b[90m│\x1b[0m\n", "NAME", "SOLVES", "SOLVED", "TYPE")
		separator := "\x1b[90m" + "│" + strings.Repeat("─", tableWidth) + "│" + "\n" + "\x1b[0m"

		sb.WriteString(header)
		sb.WriteString(separator)

		for i, k := range m.Keys {
			l := m.Levels[k]
			if m.Cursor == i {
				sb.WriteString("\x1b[90m│\x1b[0m> ")
			} else {
				sb.WriteString("\x1b[90m│\x1b[0m  ")
			}
			sb.WriteString(fmt.Sprintf("%s\t%-20s%-14s%-14s%-20s\x1b[90m│\x1b[0m\n", l.ID, strings.ToLower(l.Contract), m.solves[l.Contract], m.submissions[l.Contract], l.Type))
			if m.Cursor == i && m.descriptionShown {
				descriptionLines := strings.Split(l.Description, "\n")
				separator := "\x1b[90m" + "│" + strings.Repeat("-", tableWidth) + "│" + "\n" + "\x1b[0m"
				sb.WriteString(separator)
				for _, line := range descriptionLines {
					// Calculate the remaining space available for padding
					padding := tableWidth - len(line)

					// Ensure padding is not negative
					if padding < 0 {
						padding = 0
					}

					// Indent description lines and enclose them in the box
					sb.WriteString("\x1b[90m│" + line + strings.Repeat(" ", padding) + "│\x1b[0m\n")
				}
				// dont show seperator at the last level
				if i != len(m.Keys)-1 {
					sb.WriteString(separator)
				}
			}
		}

		sb.WriteString("\x1b[90m└" + strings.Repeat("─", tableWidth) + "┘\n\x1b[0m") // Bottom border of the box

		sb.WriteString("\n\x1b[90m↑/↓ - Navigate | ←/→ - Toggle Description | q to exit | ↩ to select \x1b[0m")

		return sb.String()
	}
}

func NewLevelList(Levels map[string]utils.Level, solves map[string]string, submissions map[string]string) (*levelListModel, error) {
	// Check terminal size
	if err := utils.CheckMinTerminalWidth(); err != nil {
		return nil, err
	}

	return &levelListModel{Levels: Levels, solves: solves, submissions: submissions}, nil
}
