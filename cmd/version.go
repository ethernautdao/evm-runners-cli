package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/utils"

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display current evm-runners version",

	RunE: func(cmd *cobra.Command, args []string) error {

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		fmt.Printf("evm-runners version %s\n", config.EVMR_VERSION)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
