package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type Submission struct {
	ID          string  `json:"id"`
	UserID      int     `json:"user_id"`
	LevelID     int     `json:"level_id"`
	Gas         float64 `json:"gas,string"`
	Size        int     `json:"size,string"`
	SubmittedAt string  `json:"submitted_at"`
	Type        string  `json:"type"`
	Username    string  `json:"user_name"`
	LevelName   string  `json:"level_name"`
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

	// Find the maximum username width
	maxUsernameWidth := 0
	for _, submission := range submissions {
		userStr := fmt.Sprintf("%s", submission.Username)
		if len(userStr) > maxUsernameWidth {
			maxUsernameWidth = len(userStr)
		}
	}

	// Add a few spaces to the width
	maxUsernameWidth += 4

	if len(submissions) > 0 {
		header := fmt.Sprintf(" #\t%-*s\t%-5s\t%s\t\t%s\n", maxUsernameWidth, "USER", strings.ToUpper(field), "DATE", "TYPE")
		separator := "\x1b[90m" + strings.Repeat("-", maxUsernameWidth+48) + "\n" + "\x1b[0m"

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

		// Convert the date string to a time.Time object and format it
		date, err := time.Parse(dateLayout, submission.SubmittedAt)
		if err != nil {
			return fmt.Sprintf("Error parsing date: %v", err)
		}
		dateStr := date.Format(displayLayout)

		if field == "gas" {
			sb.WriteString(fmt.Sprintf(" %d\t%-*s\t%d\t%s\t%s\n", i+1, maxUsernameWidth, userStr, int(submission.Gas), dateStr, submission.Type))
		} else if field == "size" {
			sb.WriteString(fmt.Sprintf(" %d\t%-*s\t%d\t%s\t%s\n", i+1, maxUsernameWidth, userStr, submission.Size, dateStr, submission.Type))
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
		default:
			return ui, tea.Quit
		}
	}
	return ui, nil
}

func (ui *CombinedLeaderboardUI) View() string {
	return "\n" + ui.GasUI.View() + "\n\n" + ui.SizeUI.View() + "\x1b[90m" + "\n\nPress any key to exit." + "\x1b[0m"
}
