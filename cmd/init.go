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
	Long: `Initializes EVM Runners by cloning the ethernautdao/evm-runners-levels.git repository into ./levels`,

	Run: func(cmd *cobra.Command, args []string) {
		subdir := "./levels"

		fmt.Println("Initializing EVM Runners ...")
		// Check if the subdirectory already exists
		if _, err := os.Stat(subdir); os.IsNotExist(err) {
			// Clone ethernautdao/evm-runners-levels.git repository
			execCmd := exec.Command("git", "clone", "git@github.com:ethernautdao/evm-runners-levels.git", subdir)
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			if err := execCmd.Run(); err != nil {
				fmt.Println("Failed to clone ethernautdao/evm-runners-levels.git:", err)
				return
			}
			fmt.Println("evm-runners-levels cloned successfully")
		} else {
			fmt.Println("Subdirectory already exists")
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}


