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

		// get level information
		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("error loading levels: %v", err)
		}

		// if argument is empty, open level list
		if len(args) == 0 {
			solves := utils.GetSolves()

			fmt.Println("Press ENTER to select a level: ")

			model := tui.NewLevelList(levels, solves)
			p := tea.NewProgram(model)

			if err := p.Start(); err != nil {
				return fmt.Errorf("error displaying level list: %v", err)
			}

			if model.Done {
				selectedLevelKey := model.Keys[model.Cursor]
				selectedLevel := model.Levels[selectedLevelKey]
				args = append(args, selectedLevel.Name)
			} else {
				return nil
			}
		}

		level := args[0]

		// check if level exists
		if _, ok := levels[level]; !ok {
			return fmt.Errorf("level %s does not exist", level)
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
				return fmt.Errorf("error displaying language selection list: %v", err)
			}

			if model.Done {
				selection = model.Options[model.Cursor]
			} else {
				return nil
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
			return fmt.Errorf("invalid language: %s", selection)
		}

		// copy level from template/src to src
		src := "./levels/template/src/" + fileToCopy
		dst := "./levels/src/" + fileToCopy

		// Check if file already exists. If yes, print warning and return
		_, err1 := os.Stat(dst)
		if err1 == nil {
			return fmt.Errorf("file already exists: %s", dst)
		}

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("error copying file: %v", err)
		}

		// copy test file from template to test
		src = "./levels/template/" + testToCopy
		dst = "./levels/test/" + testToCopy

		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("error copying file: %v", err)
		}

		fmt.Println("\nYour challenge is ready!\nCheck out the levels/src folder for your level file.\n\nGood luck!")

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
