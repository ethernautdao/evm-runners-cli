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

		// load levels
		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("error loading levels: %v", err)
		}

		// get level
		level, err := getLevel(args, config, levels)
		if err != nil {
			return err
		}
		// user aborted selection
		if level == "" {
			return nil
		}

		// get lang
		lang, err = getLang(lang)
		if err != nil {
			return err
		}
		// user aborted selection
		if lang == "" {
			return nil
		}

		testSuffixes := map[string]string{
			"sol":  "-Sol.t.sol",
			"huff": "-Huff.t.sol",
			"vy":   "-Vyper.t.sol",
		}

		if testSuffix, ok := testSuffixes[lang]; ok {
			filename := levels[level].File
			fileToCopy := filename + "." + lang
			testToCopy := filename + testSuffix

			err = copyTemplateFiles(config.EVMR_LEVELS_DIR, fileToCopy, testToCopy)
			if err != nil {
				return err
			}

		} else if lang == "no template" {
			fmt.Printf("No template file selected.\nYou can start working on your solution in '%s'!\nTo validate it, run 'evm-runners validate %s'\n", filepath.Join(config.EVMR_LEVELS_DIR, "src"), level)
			return nil
		} else {
			return fmt.Errorf("invalid language: %s", lang)
		}

		fmt.Printf("\nYour level is ready!\nOpen 'evm-runners-levels/src' to start working on it -- Good luck!\nTo validate your solution, run 'evm-runners validate %s'\n", level)

		return nil
	},
}

func getLevel(args []string, config utils.Config, levels map[string]utils.Level) (string, error) {
	// if argument is empty, open level list
	if len(args) == 0 {
		solves := utils.GetSolves(levels)

		fmt.Println("Press ENTER to select a level: ")

		// Fetch existing submission data
		var submissions map[string]string
		for key := range levels {
			sub, err := utils.FetchSubmissionData(config, "2")

			// if it errors, just return an empty string
			if err != nil {
				submissions[levels[key].Contract] = ""
			}

			if len(sub) == 0 {
				submissions[levels[key].Contract] = ""
			} else {
				submissions[levels[key].Contract] = "x"
			}

		}

		// display level list
		model := tui.NewLevelList(levels, solves, submissions)
		p := tea.NewProgram(model)

		if err := p.Start(); err != nil {
			return "", fmt.Errorf("error displaying level list: %v", err)
		}

		if model.Done {
			selectedLevelKey := model.Keys[model.Cursor]
			selectedLevel := model.Levels[selectedLevelKey]
			args = append(args, selectedLevel.Contract)

		} else {
			return "", nil
		}
	}

	level := strings.ToLower(args[0])

	// check if level exists
	if _, ok := levels[level]; !ok {
		return "", fmt.Errorf("level %s does not exist", level)
	}

	return level, nil
}

func getLang(lang string) (string, error) {
	// if lang flag is not sol, huff, or vyper => open list
	switch lang {
	case "Solidity", "solidity", "sol":
		lang = "sol"
	case "Huff", "huff":
		lang = "huff"
	case "Vyper", "vyper", "vy":
		lang = "vy"
	default:
		model := tui.NewLangListModel()
		p := tea.NewProgram(model)

		if err := p.Start(); err != nil {
			return "", fmt.Errorf("error displaying language selection list: %v", err)
		}

		if model.Done {
			lang = model.Lang[model.Cursor]
		} else {
			return "", nil
		}
	}

	return lang, nil
}

func copyTemplateFiles(levelsDir, fileToCopy, testToCopy string) error {
	// copy test file from template to test
	fmt.Printf("Copying test '%s' ...\n", testToCopy)

	src := filepath.Join(levelsDir, "template", testToCopy)
	dstTest := filepath.Join(levelsDir, "test", testToCopy)

	if err := copyFile(src, dstTest); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	fmt.Printf("Test copied successfully.\n\n")

	fmt.Printf("Copying source file '%s' ...\n", fileToCopy)

	// copy level from template/src to src
	src = filepath.Join(levelsDir, "template", "src", fileToCopy)
	dstSource := filepath.Join(levelsDir, "src", fileToCopy)

	// Check if file already exists. If yes, ask if overwrite is wanted
	_, err := os.Stat(dstSource)
	if !os.IsNotExist(err) {
		fmt.Printf("Source file already exists in '%s'. Overwrite? (y/n): ", dstSource)
		var overwrite string
		_, err := fmt.Scanln(&overwrite)
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		// print new line
		fmt.Printf("\n")

		if overwrite != "y" && overwrite != "Y" {
			fmt.Printf("Not overwriting '%s'\n", fileToCopy)
			return nil
		}
	}

	if err := copyFile(src, dstSource); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	return nil
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
