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
	Run: func(cmd *cobra.Command, args []string) {
		levels, err := config.LoadLevels()

		if err != nil {
			fmt.Println(err)
		}

		model := tui.NewLevelList(levels)
		p := tea.NewProgram(model)

		if err := p.Start(); err != nil {
			fmt.Println("Error creating level list:", err)
			return
		}

		if model.Done {
			selectedLevelKey := model.Keys[model.Cursor]
			selectedLevel := model.Levels[selectedLevelKey]
			// Use selectedLevel for your needs
			fmt.Printf("Selected level: %v\n", selectedLevel)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
