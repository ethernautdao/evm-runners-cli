package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// aboutCmd represents the about command
var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Information about evm-runners",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`                                                                 
  _____   ___ __ ___        _ __ _   _ _ __  _ __   ___ _ __ ___ 
 / _ \ \ / / '_ ' _ \ _____| '__| | | | '_ \| '_ \ / _ \ '__/ __|
|  __/\ V /| | | | | |_____| |  | |_| | | | | | | |  __/ |  \__ \
 \___| \_/ |_| |_| |_|     |_|   \__,_|_| |_|_| |_|\___|_|  |___/
                                                                 `)

		fmt.Println("A terminal-based game for developers with EVM-based challenges")
		fmt.Println("\nSponsored by @EthernautDAO")
		fmt.Println("Authors: @0xkarmacoma, @beskay0x, @kyre_rs")
		fmt.Println("Website: evmr.sh")

		fmt.Println("\nevm-runners is not a regular CTF:")
		fmt.Println(" - Score based, not hack based")
		fmt.Println(" - Dual scores (gas and codesize), not a compound score")
		fmt.Println(" - No time limit")
		fmt.Println(" - Linear progression: Challenges get more complex as you progress")
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
