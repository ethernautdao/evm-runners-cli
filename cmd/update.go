package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/utils"

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update evm-runners levels",
	RunE: func(cmd *cobra.Command, args []string) error {
		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		fmt.Printf("Updating evm-runners levels...\n\n")
		if _, err := os.Stat(config.EVMR_LEVELS_DIR); os.IsNotExist(err) {
			return fmt.Errorf("evm-runners directory not found, run 'evmr init' first!\n")
		} else {
			execCmd := exec.Command("git", "pull")
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			execCmd.Dir = config.EVMR_LEVELS_DIR
			if err := execCmd.Run(); err != nil {
				return fmt.Errorf("error updating the level directory: %v", err)
			}
		}
		fmt.Println("\nDone!")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
