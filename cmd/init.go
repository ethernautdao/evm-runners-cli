package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes EVM Runners",

	Run: func(cmd *cobra.Command, args []string) {
		subdir := "./levels"

		fmt.Println("Initializing EVM Runners ...")
		// Check if the subdirectory already exists
		if _, err := os.Stat(subdir); os.IsNotExist(err) {
			// Clone the subdirectory from Github
			cmd := exec.Command("git", "clone", "git@github.com:ethernautdao/evm-runners-levels.git", subdir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Println("Failed to clone subdirectory:", err)
				return
			}
			fmt.Println("Subdirectory cloned successfully")
		} else {
			fmt.Println("Subdirectory already exists")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}


