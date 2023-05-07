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
		fmt.Println("Website: https://evmr.sh")
		fmt.Println("Discord: https://discord.gg/RQ5WYDxUF3")

		fmt.Println("\nevm-runners is not just a regular CTF:")
		fmt.Println(" - Gradual introduction to the EVM")
		fmt.Println(" - Score based, not hack based")
		fmt.Println(" - Dual scores (gas and codesize), not a compound score")
		fmt.Println(" - Linear progression: Challenges get more complex as you progress")
		fmt.Println(" - Multiplayer: Compare your scores with other players")
		fmt.Println(" - Singleplayer: Play at your own pace, learn usefull skills")
		fmt.Println(" - Work on solutions in any EVM-language (Solidity, Vyper, Huff, etc.)")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
