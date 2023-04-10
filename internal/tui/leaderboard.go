package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Submission struct {
	ID       string  `json:"id"`
	LevelID  int     `json:"level_id"`
	UserID   int     `json:"user_id"`
	Bytecode string  `json:"bytecode"`
	Gas      float64 `json:"gas,string"` // Use string tag to handle the gas value as string
	Size     int     `json:"size,string"`
}

type LeaderboardUI struct {
	content string
}

type CombinedLeaderboardUI struct {
	GasUI  *LeaderboardUI
	SizeUI *LeaderboardUI
}

func NewLeaderboardUI(submissions []Submission, field string) *LeaderboardUI {
	ui := &LeaderboardUI{
		content: leaderboardTable(submissions, field),
	}
	return ui
}

func leaderboardTable(submissions []Submission, field string) string {
	var sb strings.Builder

	header := fmt.Sprintf("USER ID | %s\n", strings.ToUpper(field))
	headerSeparator := strings.Repeat("-", len(header)-1) + "\n"

	sb.WriteString(header)
	sb.WriteString(headerSeparator)

	for _, submission := range submissions {
		if field == "gas" {
			sb.WriteString(fmt.Sprintf("%-7d | %.2f\n", submission.UserID, submission.Gas))
		} else if field == "size" {
			sb.WriteString(fmt.Sprintf("%-7d | %d\n", submission.UserID, submission.Size))
		}
	}

	return sb.String()
}

func (ui *LeaderboardUI) View() string {
	return ui.content
}

func (ui *CombinedLeaderboardUI) Init() tea.Cmd {
	return nil
}

func (ui *CombinedLeaderboardUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return ui, tea.Quit
		}
	}
	return ui, nil
}

func (ui *CombinedLeaderboardUI) View() string {
	return ui.GasUI.View() + "\n\n" + ui.SizeUI.View()
}
