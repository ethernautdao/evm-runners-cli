package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
)

type LeaderboardUI struct {
	content string
}

type CombinedLeaderboardUI struct {
	GasUI  *LeaderboardUI
	SizeUI *LeaderboardUI
}

func NewLeaderboardUI(submissions []utils.SubmissionData, field string) (*LeaderboardUI, error) {
	// Check terminal size
	if err := utils.CheckMinTerminalWidth(); err != nil {
		return nil, err
	}

	ui := &LeaderboardUI{
		content: leaderboardTable(submissions, field),
	}
	return ui, nil
}

func leaderboardTable(submissions []utils.SubmissionData, field string) string {
	var sb strings.Builder

	tableWidth := 75

	if len(submissions) > 0 {

		var headlineText string
		if field == "gas" {
			headlineText = "GAS LEADERBOARD"
		} else if field == "size" {
			headlineText = "SIZE LEADERBOARD"
		}

		headline := "\x1b[1m" + headlineText + "\x1b[0m" + "\n\n"
		header := fmt.Sprintf("\x1b[90m│\x1b[0m #\t%-22s%-14s%-20s%-12s\x1b[90m│\x1b[0m\n", "USER", strings.ToUpper(field), "DATE", "TYPE")
		separator := "\x1b[90m" + "│" + strings.Repeat("─", tableWidth) + "│" + "\n" + "\x1b[0m"

		// Calculate padding for the headline
		headlinePadding := (tableWidth - len(headlineText)) / 2
		padding := strings.Repeat(" ", headlinePadding)
		centeredHeadline := padding + headline

		sb.WriteString(centeredHeadline)
		sb.WriteString("\x1b[90m┌" + strings.Repeat("─", tableWidth) + "┐\n\x1b[0m") // Top border of the box
		sb.WriteString(header)
		sb.WriteString(separator)
	} else {
		// No submissions, display a message
		sb.WriteString(fmt.Sprintf("No submissions available for the %s leaderboard!", field))
		return sb.String()
	}

	dateLayout := "2006-01-02T15:04:05.000Z"
	displayLayout := "Jan 02 2006"

	for i, submission := range submissions {
		userStr := fmt.Sprintf("%s", submission.Username)

		// Check if the userStr is longer than 20 characters
		if len(userStr) > 20 {
			// Truncate the userStr to 20 characters
			userStr = userStr[:20]

			// Replace the last two characters with ".."
			userStr = userStr[:18] + ".."
		}

		// Convert the date string to a time.Time object and format it
		date, err := time.Parse(dateLayout, submission.SubmittedAt)
		if err != nil {
			return fmt.Sprintf("Error parsing date: %v", err)
		}
		dateStr := date.Format(displayLayout)

		if field == "gas" {
			sb.WriteString(fmt.Sprintf("\x1b[90m│\x1b[0m %d\t%-22s%-14s%-20s%-12s\x1b[90m│\x1b[0m\n", i+1, userStr, submission.Gas, dateStr, submission.Type))
		} else if field == "size" {
			sb.WriteString(fmt.Sprintf("\x1b[90m│\x1b[0m %d\t%-22s%-14s%-20s%-12s\x1b[90m│\x1b[0m\n", i+1, userStr, submission.Size, dateStr, submission.Type))
		}
	}

	sb.WriteString("\x1b[90m└" + strings.Repeat("─", tableWidth) + "┘\n\x1b[0m") // Bottom border of the box

	return sb.String()
}

func (ui *LeaderboardUI) View() string {
	return ui.content
}

func (ui *CombinedLeaderboardUI) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (ui *CombinedLeaderboardUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		default:
			return ui, tea.Quit
		}
	}
	return ui, nil
}

func (ui *CombinedLeaderboardUI) View() string {
	return "\n" + ui.GasUI.View() + "\n\n" + ui.SizeUI.View() + "\x1b[90m" + "\n\nPress any key to exit." + "\x1b[0m"
}
