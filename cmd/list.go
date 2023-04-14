package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/tui"
import "github.com/ethernautdao/evm-runners-cli/internal/config"

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all emv-runners levels",
	RunE: func(cmd *cobra.Command, args []string) error {
		levels, err := config.LoadLevels()
		if err != nil {
			fmt.Println("Error loading levels")
			return err
		}

		solves := config.GetSolves()

		model := tui.NewLevelList(levels, solves)
		p := tea.NewProgram(model)

		if err := p.Start(); err != nil {
			fmt.Println("Error displaying level list")
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
