package tui

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
)

type model struct {
    options []string
    selected int
}

var choice string = ""

func GetChoice() string {
	return choice
}

func initialModel() model {
	return model{
		options:  []string{"Solidity", "Huff"},
        selected: 0,
	}
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
        case "up":
            m.selected--
            if m.selected < 0 {
                m.selected = len(m.options) - 1
            }
        case "down":
            m.selected++
            if m.selected >= len(m.options) {
                m.selected = 0
            }
        case "enter":
			choice = m.options[m.selected]
			return m, tea.Quit
		}
    }

    return m, nil
}

func (m model) View() string {
    s := "Select a language you want to choose:\n\n"

    for i, option := range m.options {
        // Add a ">" symbol before the selected option
        if i == m.selected {
            s += "> "
        } else {
            s += "  "
        }
        s += option + "\n"
    }

    return s
}

func RunBubbleTea() string {

	m := initialModel()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
        os.Exit(1)
	}

	return m.options[m.selected]
}