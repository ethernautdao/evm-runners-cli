package tui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"time"
)

type Submission struct {
	ID            string  `json:"id"`
	UserID        int     `json:"user_id"`
	LevelID       int     `json:"level_id"`
	Gas           float64 `json:"gas,string"`
	Size          int     `json:"size,string"`
	SubmittedAt   string  `json:"submitted_at"`
	Username      string  `json:"user_name"`
	Discriminator int     `json:"discriminator"`
	LevelName     string  `json:"level_name"`
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
		userStr := fmt.Sprintf("%s#%04d", submission.Username, submission.Discriminator)
		if len(userStr) > maxUsernameWidth {
			maxUsernameWidth = len(userStr)
		}
	}

	// Add a few spaces to the width
	maxUsernameWidth += 4

	if len(submissions) > 0 {
		/*         levelName := submissions[0].LevelName
		           title := fmt.Sprintf("%s Leaderboard for level %s", strings.ToUpper(field), levelName) */

		header := fmt.Sprintf(" #\t%-*s\t%-5s\t%s\n", maxUsernameWidth, "User", strings.ToUpper(field), "Date")
		separator := strings.Repeat("-", maxUsernameWidth+38) + "\n"

		// Calculate padding for centered title
		/*         tableWidth := maxUsernameWidth + 38
		           titlePadding := (tableWidth - len(title)) / 2
		           titleLine := fmt.Sprintf("%-*s%s\n", titlePadding, "", title) */

		//sb.WriteString(titleLine)
		//sb.WriteString("\n")
		sb.WriteString(header)
		sb.WriteString(separator)
	} else {
		// No submissions, display a message
		sb.WriteString("No submissions available for this leaderboard.")
		return sb.String()
	}

	dateLayout := "2006-01-02T15:04:05.000Z"
	displayLayout := "Jan 02 2006"

	for i, submission := range submissions {
		userStr := fmt.Sprintf("%s#%04d", submission.Username, submission.Discriminator)

		// Convert the date string to a time.Time object and format it
		date, err := time.Parse(dateLayout, submission.SubmittedAt)
		if err != nil {
			return fmt.Sprintf("Error parsing date: %v", err)
		}
		dateStr := date.Format(displayLayout)

		if field == "gas" {
			sb.WriteString(fmt.Sprintf(" %d\t%-*s\t%d\t%s\n", i+1, maxUsernameWidth, userStr, int(submission.Gas), dateStr))
		} else if field == "size" {
			sb.WriteString(fmt.Sprintf(" %d\t%-*s\t%d\t%s\n", i+1, maxUsernameWidth, userStr, submission.Size, dateStr))
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
	return "\n" + ui.GasUI.View() + "\n\n" + ui.SizeUI.View() + "\n\nPress any key to exit."
}
