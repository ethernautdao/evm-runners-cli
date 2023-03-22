package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/tui"

import (
	"fmt"
    "io/ioutil"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Starts a challenge",
    Long:  `Starts a challenge by selecting a level with the --level flag and choosing a language`,
    
	Run: func(cmd *cobra.Command, args []string) {
    level, _ := cmd.Flags().GetString("level")

	var extension string

	// display list
	tui.RunBubbleTea()
	
	// get user choice
	switch tui.GetChoice() {
	case "Solidity":
		extension = ".sol"
	case "Huff":
		extension = ".huff"
	default:
		extension = ".sol"
	} 

	src := "./levels/src/template/" + level + extension
	dst := "./levels/src/" + level + extension

	if err := copyFile(src, dst); err != nil {
		fmt.Printf("There's been an error: %v", err)
	}

	fmt.Println("\nYour challenge is ready! Check out the levels/src folder for your level file. Good luck!")

    },
}

func copyFile(src, dst string) error {
    input, err := ioutil.ReadFile(src)
    if err != nil {
        return err
    }
    if err := ioutil.WriteFile(dst, input, 0644); err != nil {
        return err
    }
    return nil
}

func init() {
    rootCmd.AddCommand(startCmd)
    startCmd.Flags().StringP("level", "l", "", "Select a level")
    startCmd.MarkFlagRequired("level")
}
