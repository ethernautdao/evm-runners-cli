package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/tui"
import "github.com/ethernautdao/evm-runners-cli/internal/utils"

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
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
		src := filepath.Join(".", "levels", "template", "src", fileToCopy)
		dstSource := filepath.Join(".", "levels", "src", fileToCopy)

		// add blank line
		fmt.Printf("\n")

		// Check if file already exists. If yes, ask if overwrite is wanted
		_, err = os.Stat(dstSource)
		if !os.IsNotExist(err) {
			fmt.Printf("File %s already exists. Overwrite? (y/n): ", dstSource)
			var overwrite string
			_, err := fmt.Scanln(&overwrite)
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}
			if overwrite != "y" && overwrite != "Y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		if err := copyFile(src, dstSource); err != nil {
			return fmt.Errorf("error copying file: %v", err)
		}

		// copy test file from template to test
		src = filepath.Join(".", "levels", "template", testToCopy)
		dstTest := filepath.Join(".", "levels", "test", testToCopy)

		if err := copyFile(src, dstTest); err != nil {
			return fmt.Errorf("error copying file: %v", err)
		}

		fmt.Printf("Your challenge is ready!\nCheck out %s for your level file and %s for your test file.\n\nGood luck!\n", dstSource, dstTest)

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
