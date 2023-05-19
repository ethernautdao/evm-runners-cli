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
	Short: "Lists all evm-runners levels",

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

		// Fetch existing submission data
		submissions := make(map[string]string) // Initialize the submissions map
		for key := range levels {
			sub, err := utils.FetchSubmissionData(config, levels[key].ID)

			// if it errors, just return an empty string
			if err != nil {
				submissions[levels[key].Contract] = ""
			}

			if len(sub) == 0 {
				submissions[levels[key].Contract] = ""
			} else {
				submissions[levels[key].Contract] = "x"
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
