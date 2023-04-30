package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/tui"
import "github.com/ethernautdao/evm-runners-cli/internal/utils"

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all evm-runners levels",

	RunE: func(cmd *cobra.Command, args []string) error {
		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("error loading levels: %v", err)
		}

		// TODO: Add timeout if no resp

		solves := utils.GetSolves()

		model := tui.NewLevelList(levels, solves)
		p := tea.NewProgram(model)

		if err := p.Start(); err != nil {
			return fmt.Errorf("error displaying level list: %v", err)
		}

		// TODO: Show last solve time?

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
