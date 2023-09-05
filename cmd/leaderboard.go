package cmd

import (
	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethernautdao/evm-runners-cli/internal/tui"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

var leaderboardCmd = &cobra.Command{
	Use:   "leaderboard [level]",
	Short: "Display the gas and codesize leaderboard for a specific level",

	RunE: func(cmd *cobra.Command, args []string) error {
		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		// load levels
		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("error loading levels: %v", err)
		}

		// get level
		level, err := GetLevel(args, config, levels)
		if err != nil {
			return err
		}
		// user aborted selection
		if level == "" {
			return nil
		}

		return displayLeaderboard(levels[level].ID)
	},
}

func fetchLeaderboardData(url string) ([]utils.SubmissionData, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var leaderboardData []utils.SubmissionData
	if err := json.Unmarshal(body, &leaderboardData); err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	// Limit the array length to 10 if it exceeds
	if len(leaderboardData) > 10 {
		leaderboardData = leaderboardData[:10]
	}

	return leaderboardData, nil
}

func displayLeaderboard(levelId string) error {
	config, err := utils.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	gasURL := fmt.Sprintf("%ssubmissions/leaderboard/gas/%s", config.EVMR_SERVER, levelId)
	sizeURL := fmt.Sprintf("%ssubmissions/leaderboard/size/%s", config.EVMR_SERVER, levelId)

	// Fetch gas leaderboard data
	gasLeaderboardData, err := fetchLeaderboardData(gasURL)
	if err != nil {
		return fmt.Errorf("error fetching gas leaderboard data: %v", err)
	}

	// Fetch size leaderboard data
	sizeLeaderboardData, err := fetchLeaderboardData(sizeURL)
	if err != nil {
		return fmt.Errorf("error fetching size leaderboard data: %v", err)
	}

	// Initialize the BubbleTea UI
	gasUI, err := tui.NewLeaderboardUI(gasLeaderboardData, "gas")
	if err != nil {
		// can only fail if terminal width is < required width, so we bubble up error
		return err
	}
	sizeUI, err := tui.NewLeaderboardUI(sizeLeaderboardData, "size")
	if err != nil {
		// can only fail if terminal width is < required width, so we bubble up error
		return err
	}

	// Combine the views of gasUI and sizeUI
	m := &tui.CombinedLeaderboardUI{GasUI: gasUI, SizeUI: sizeUI}

	// Run the BubbleTea program
	err = tea.NewProgram(m).Start()
	if err != nil {
		return fmt.Errorf("error running BubbleTea UI: %v", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(leaderboardCmd)
}
