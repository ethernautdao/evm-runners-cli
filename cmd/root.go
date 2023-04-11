package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "evmrunners",
	Short: "A terminal-based game for developers with EVM-based challenges",
	Long: `A terminal-based game for developers with EVM-based challenges.

How to play:
1. Run 'evmrunners init' to initialize the game
2. Run 'evmrunners list' to list all available levels
3. Run 'evmrunners start <level>' to start solving a level
4. Run 'evmrunners validate <level>' to validate your solution
5. Run 'evmrunners submit <level>' to submit your solution`,
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
