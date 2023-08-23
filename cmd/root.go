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

How to start playing:
1. 'evmr init' to initialize evm-runners
2. 'evmr start' to start solving a level
3. 'evmr validate <level>' to test your solution

How to submit a solution:
1. 'evmr auth' to register your account
2. 'evmr submit <level>' to submit your solution
3. 'evmr leaderboard <level>' to see your position on the leaderboard

For a list of all available levels, use 'evmr levels'`,
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
