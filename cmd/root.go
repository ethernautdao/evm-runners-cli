package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "evm-runners",
	Short: "A terminal-based game for developers with EVM-based challenges",
	Long: `A terminal-based game for developers with EVM-based challenges.

How to play:
1. Run 'evm-runners init' to initialize the game
2. Run 'evm-runners list' to list all available levels
3. Run 'evm-runners start <level>' to start solving a level
4. Run 'evm-runners validate <level>' to validate your solution
5. Run 'evm-runners submit <level>' to submit your solution`,
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
