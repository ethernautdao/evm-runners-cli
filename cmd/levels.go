package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/tui"
import "github.com/ethernautdao/evm-runners-cli/internal/utils"

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// levelsCmd represents the levels command
var levelsCmd = &cobra.Command{
	Use:   "levels",
	Short: "List all available evm-runners levels",

	RunE: func(cmd *cobra.Command, args []string) error {

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("error loading levels: %v", err)
		}

		// get amount of solves for each level
		solves := utils.GetSolves(levels)

		// Initialize the submissions map
		submissions := make(map[string]string)
		for key := range levels {
			submissions[levels[key].Contract] = ""
		}

		// Fetch existing submission data if user authenticated
		if config.EVMR_TOKEN != "" {
			// we explicitly ignore checking the error here
			sub, _ := utils.FetchSubmissionData(config)

			for _, item := range sub {
				submissions[item.LevelName] = "x"
			}
		}

		// display level list
		model := tui.NewLevelList(levels, solves, submissions)
		p := tea.NewProgram(model)

		if err := p.Start(); err != nil {
			return fmt.Errorf("error displaying level list: %v", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(levelsCmd)
}
