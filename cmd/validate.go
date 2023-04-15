package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ethernautdao/evm-runners-cli/internal/utils"
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

		// get filename and test contract of level
		filename := levels[level].FileName
		testContract := levels[level].TestContract

		// if bytecode is provided, set the BYTECODE env variable
		if bytecode != "" {
			bytecode = utils.CheckValidBytecode(bytecode)
			os.Setenv("BYTECODE", bytecode)
		} else {
			solutionType := utils.CheckSolutionFile(filename, lang)
			if solutionType == "nil" {
				return nil
			}

			// choose specific test contract (either sol or huff version)
			if solutionType == "sol" {
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
