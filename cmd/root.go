package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "evm-runners",
	Short: "A terminal-based game for developers with EVM-based levels",
	Long: `A terminal-based game for developers with EVM-based levels.

How to play:
1. 'evm-runners init' to initialize the game
2. 'evm-runners auth' to authenticate with the server
3. 'evm-runners start' to start solving a level
4. 'evm-runners validate <level>' to validate your solution
5. 'evm-runners submit <level>' to submit your solution and get placed on the leaderboard

For a list of all available levels, use 'evm-runners levels'`,
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
