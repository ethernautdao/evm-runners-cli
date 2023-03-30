package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/config"

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all emv-runners levels",

	Run: func(cmd *cobra.Command, args []string) {
		levels, err := config.LoadLevels();

		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(levels)

		fmt.Println(levels[0].FileName)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
