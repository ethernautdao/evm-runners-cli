package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ethernautdao/evm-runners-cli/internal/config"
	"github.com/spf13/cobra"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate <level>",
	Short: "Validates a level",
	Long: `Validates a level by running the predefined Foundry tests against 
the submitted solution file (either .huff or .sol) or against the provided bytecode, if set.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		bytecode, _ := cmd.Flags().GetString("bytecode")
		lang, _ := cmd.Flags().GetString("lang")

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

		// get filename and test contract of level
		filename := levels[level].FileName
		testContract := levels[level].TestContract

		// if bytecode is provided, set the BYTECODE env variable
		if bytecode != "" {
			os.Setenv("BYTECODE", bytecode)
		} else if lang != "" {
			// if language flag is set, check if the corresponding solution file exists
			if lang == "sol" {
				_, err := os.Stat(fmt.Sprintf("./levels/src/%s.sol", filename))
				if os.IsNotExist(err) {
					fmt.Println("No Solidity solution file found. Add a solution file or submit bytecode with the --bytecode flag!")
					return nil
				}
				testContract = testContract + "Sol"
			} else if lang == "huff" {
				_, err := os.Stat(fmt.Sprintf("./levels/src/%s.huff", filename))
				if os.IsNotExist(err) {
					fmt.Println("No Huff solution file found. Add a solution file or submit bytecode with the --bytecode flag!")
					return nil
				}
				testContract = testContract + "Huff"
			} else {
				fmt.Println("Invalid language flag. Please use either 'sol' or 'huff'.")
				return nil
			}
		} else {
			// Check existence of solution files if no bytecode flag and no language flag is set
			_, err1 := os.Stat(fmt.Sprintf("./levels/src/%s.sol", filename))
			_, err2 := os.Stat(fmt.Sprintf("./levels/src/%s.huff", filename))

			if os.IsNotExist(err1) && os.IsNotExist(err2) {
				fmt.Println("No solution file found. Add a solution file or submit bytecode with the --bytecode flag!")
				return nil
			} else if err1 == nil && err2 == nil {
				fmt.Println("More than one solution file found!\nDelete a solution file or use the --language flag to choose which one to validate.")
				return nil
			}

			// choose specific test contract (either sol or huff version)
			if err1 == nil {
				testContract = testContract + "Sol"
			} else {
				testContract = testContract + "Huff"
			}
		}

		// TODO: check if test files got tampered with

		// Create the command to be run in the subdirectory
		execCmd := exec.Command("forge", "test", "--match-contract", testContract, "-vv")

		// Set the working directory to the subdirectory
		execCmd.Dir = "./levels/"

		// Capture the standard output and standard error of the command
		output, err := execCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s", output)
		}

		// Print the output of the command to the console
		fmt.Println(string(output))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringP("bytecode", "b", "", "The creation bytecode to submit")
	validateCmd.Flags().StringP("lang", "l", "", "The lang of the solution file. Either 'sol' or 'huff'")
}
