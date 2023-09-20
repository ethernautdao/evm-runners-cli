package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "evmr",
	Short: "A terminal-based game for developers with EVM-based levels",
	Long: `A terminal-based game for developers with EVM-based levels.

How to play:
  1. 'evmr init' - Initialize evm-runners.
  2. 'evmr start [level]' - Begin solving a level.
  3. 'evmr validate <level>' - Test your solution.

How to submit a solution:
  1. 'evmr auth discord' - Register your account.
  2. 'evmr submit <level>' - Submit your solution.
  3. 'evmr leaderboard [level]' - Check your leaderboard position.

For a list of all available levels, use 'evmr levels'.

Note: Arguments in <> are required, while arguments in [] are optional.`,
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
