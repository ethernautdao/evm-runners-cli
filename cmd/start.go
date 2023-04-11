package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/tui"
import "github.com/ethernautdao/evm-runners-cli/internal/config"

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"io/ioutil"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <level>",
	Short: "Starts a challenge",
	Long:  `Starts a challenge by copying the respective template files from levels/template to levels/src and levels/test`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Please provide a level\n")
		}
		level := args[0]

		// get level information
		levels, err := config.LoadLevels()
		if err != nil {
			fmt.Println("Error loading levels")
			return err
		}

		// check if level exists
		if _, ok := levels[level]; !ok {
			fmt.Println("Invalid level")
			return nil
		}

		// get filename of level and declare test file
		filename := levels[level].FileName

		var fileToCopy string
		var testToCopy string

		model := tui.NewLangListModel()
		p := tea.NewProgram(model)

		if err := p.Start(); err != nil {
			fmt.Println("Error displaying language selection")
			return err
		}

		var selection string

		if model.Done {
			selection = model.Options[model.Cursor]
			//fmt.Printf("Selected level: %v\n", model.Options[model.Cursor])
		}

		switch selection {
		case "Solidity":
			fileToCopy = filename + ".sol"
			testToCopy = filename + "-Sol.t.sol"
		case "Huff":
			fileToCopy = filename + ".huff"
			testToCopy = filename + "-Huff.t.sol"
		}

		// TODO: Check if file already exists. If yes, print warning and return

		// copy level from template/src to src
		src := "./levels/template/src/" + fileToCopy
		dst := "./levels/src/" + fileToCopy

		if err := copyFile(src, dst); err != nil {
			fmt.Printf("Error copying file")
			return err
		}

		// copy test file from template to test
		src = "./levels/template/" + testToCopy
		dst = "./levels/test/" + testToCopy

		if err := copyFile(src, dst); err != nil {
			fmt.Printf("Error copying file")
			return err
		}

		fmt.Println("\nYour challenge is ready! Check out the levels/src folder for your level file. Good luck!")

		return nil
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
}
