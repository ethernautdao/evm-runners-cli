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
		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("error loading levels: %v", err)
		}

		// get amount of solves for each level
		solves := utils.GetSolves(levels)

		// display level list
		model := tui.NewLevelList(levels, solves)
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
