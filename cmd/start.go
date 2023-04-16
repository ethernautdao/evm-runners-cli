package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/tui"
import "github.com/ethernautdao/evm-runners-cli/internal/utils"

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <level>",
	Short: "Starts a challenge",
	Long:  `Starts a challenge by copying the respective template files from levels/template to levels/src and levels/test`,

	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")

		if len(args) == 0 {
			return fmt.Errorf("Please provide a level\n")
		}
		level := args[0]

		// get level information
		levels, err := utils.LoadLevels()
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
		filename := levels[level].File

		var fileToCopy string
		var testToCopy string

		var selection string

		// if lang flag is not sol or huff, open list
		if lang != "sol" && lang != "huff" {
			model := tui.NewLangListModel()
			p := tea.NewProgram(model)

			if err := p.Start(); err != nil {
				fmt.Println("Error displaying language selection")
				return err
			}

			if model.Done {
				selection = model.Options[model.Cursor]
			}
		} else {
			selection = lang
		}

		switch selection {
		case "Solidity", "sol", "solidity":
			fileToCopy = filename + ".sol"
			testToCopy = filename + "-Sol.t.sol"
		case "Huff", "huff":
			fileToCopy = filename + ".huff"
			testToCopy = filename + "-Huff.t.sol"
		default:
			fmt.Println("Invalid language")
			return nil
		}

		// copy level from template/src to src
		src := "./levels/template/src/" + fileToCopy
		dst := "./levels/src/" + fileToCopy

		// Check if file already exists. If yes, print warning and return
		_, err1 := os.Stat(dst)
		if err1 == nil {
			fmt.Println("File already exists. Please delete it first.")
			return nil
		}

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

	startCmd.Flags().StringP("lang", "l", "", "The language you want to choose. Either 'sol' or 'huff'")
}
