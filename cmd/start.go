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
	"strings"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <level>",
	Short: "Starts a level",
	Long: `Starts a level by copying the respective template files from 
evm-runners-levels/template to evm-runners-levels/src and evm-runners-levels/test`,

	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

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

		level := strings.ToLower(args[0])

		// check if level exists
		if _, ok := levels[level]; !ok {
			return fmt.Errorf("level %s does not exist", level)
		}

		// if lang flag is not sol, huff, or vyper => open list
		switch lang {
		case "Solidity", "solidity", "sol":
			lang = "sol"
		case "Huff", "huff":
			lang = "huff"
		case "Vyper", "vyper", "vy":
			lang = "vyper"
		default:
			model := tui.NewLangListModel()
			p := tea.NewProgram(model)

			if err := p.Start(); err != nil {
				return fmt.Errorf("error displaying language selection list: %v", err)
			}

			if model.Done {
				lang = model.Options[model.Cursor]
			} else {
				return nil
			}
		}

		// get filename of level and declare test file
		filename := levels[level].File

		var fileToCopy string
		var testToCopy string

		switch lang {
		case "Solidity", "sol", "solidity":
			fileToCopy = filename + ".sol"
			testToCopy = filename + "-Sol.t.sol"
		case "Huff", "huff":
			fileToCopy = filename + ".huff"
			testToCopy = filename + "-Huff.t.sol"
		case "Vyper", "vyper":
			fileToCopy = filename + ".vy"
			testToCopy = filename + "-Vyper.t.sol"
		case "no template":
			fmt.Printf("No template file choosen.\nYou can start working on your solution in '%s'!\nTo validate it, run 'evm-runners validate <level>'\n", filepath.Join(config.EVMR_LEVELS_DIR, "src"))
			return nil
		default:
			return fmt.Errorf("invalid language: %s", lang)
		}

		// copy level from template/src to src
		src := filepath.Join(config.EVMR_LEVELS_DIR, "template", "src", fileToCopy)
		dstSource := filepath.Join(config.EVMR_LEVELS_DIR, "src", fileToCopy)

		// Check if file already exists. If yes, ask if overwrite is wanted
		_, err = os.Stat(dstSource)
		if !os.IsNotExist(err) {
			fmt.Printf("File %s already exists in evm-runners-levels/src/. Overwrite? (y/n): ", fileToCopy)
			var overwrite string
			_, err := fmt.Scanln(&overwrite)
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}

			// print new line
			fmt.Printf("\n")

			if overwrite != "y" && overwrite != "Y" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		if err := copyFile(src, dstSource); err != nil {
			return fmt.Errorf("error copying file: %v", err)
		}

		// copy test file from template to test
		src = filepath.Join(config.EVMR_LEVELS_DIR, "template", testToCopy)
		dstTest := filepath.Join(config.EVMR_LEVELS_DIR, "test", testToCopy)

		if err := copyFile(src, dstTest); err != nil {
			return fmt.Errorf("error copying file: %v", err)
		}

		fmt.Printf("Your level is ready!\nOpen evm-runners-levels/src to start working on it -- Good luck!\nTo validate your solution, run 'evm-runners validate <level>'\n")

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
