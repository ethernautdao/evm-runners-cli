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
		fmt.Println("\nSponsored by \x1b[94m@EthernautDAO\x1b[0m")
		fmt.Println("Authors: \x1b[94m@0xkarmacoma\x1b[0m, \x1b[94m@beskay0x\x1b[0m, \x1b[94m@kyre_rs\x1b[0m")
		fmt.Println("")
		fmt.Println("Website: \x1b[94mhttps://evmr.sh\x1b[0m")
		fmt.Println("Discord: \x1b[94mhttps://discord.gg/2TwURWvnVT\x1b[0m")

		fmt.Println("\nevm-runners is not just a regular CTF:")
		fmt.Println(" - Gradual introduction to the EVM")
		fmt.Println(" - Score based, not hack based")
		fmt.Println(" - Dual scores (gas and codesize), not a compound score")
		fmt.Println(" - Linear progression: Challenges get more complex as you progress")
		fmt.Println(" - Multiplayer: Compare your scores with other players")
		fmt.Println(" - Singleplayer: Play at your own pace, learn usefull skills")
		fmt.Println(" - Work on solutions in a language of your chocie (Solidity, Vyper, Yul, Huff)")
		fmt.Println("")
	},
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}
