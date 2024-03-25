package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "evmr",
	Short: "A terminal-based game for developers with EVM-based levels",
	Long: `                                                            
  _____   ___ __ ___        _ __ _   _ _ __  _ __   ___ _ __ ___ 
 / _ \ \ / / '_ ' _ \ _____| '__| | | | '_ \| '_ \ / _ \ '__/ __|
|  __/\ V /| | | | | |_____| |  | |_| | | | | | | |  __/ |  \__ \
 \___| \_/ |_| |_| |_|     |_|   \__,_|_| |_|_| |_|\___|_|  |___/
                                                                 	
A terminal-based game for developers with EVM-based levels.

How to play:
  1. 'evmr init' - Initialize evm-runners.
  2. 'evmr start [level]' - Begin solving a level.
  3. 'evmr validate <level>' - Test your solution.
  4. 'evmr submit <level>' - Submit your solution.

Arguments in <> are required, while arguments in [] are optional.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
