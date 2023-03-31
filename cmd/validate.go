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
	Use:   "validate",
	Short: "Validates a level",
	Long:  `Validates a level by running the predefined Foundry tests against the submitted solution file (either .huff or .sol) or against the provided bytecode, if set.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		level, _ := cmd.Flags().GetString("level")
		bytecode, _ := cmd.Flags().GetString("bytecode")

		// get level information
		levels, err := config.LoadLevels()
		if err != nil {
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

		// Check if the level is valid
		if testContract == "" {
			fmt.Println("Invalid level")
			return nil
		}

		// if bytecode is provided, set the BYTECODE env variable
		if bytecode != "" {
			os.Setenv("BYTECODE", bytecode)
		} else {
			// Check existence of solution files if no bytecode is provided
			_, err1 := os.Stat(fmt.Sprintf("./levels/src/%s.sol", filename))
			_, err2 := os.Stat(fmt.Sprintf("./levels/src/%s.huff", filename))

			if os.IsNotExist(err1) && os.IsNotExist(err2) {
				fmt.Println("No solution file found. Add a solution file or submit bytecode with the --bytecode flag!")
				return nil
			} else if err1 == nil && err2 == nil {
				fmt.Println("More than one solution file found. Delete the one you dont want to validate!")
				return nil
			}
		}

		// Create the command to be run in the subdirectory
		execCmd := exec.Command("forge", "test", "--match-contract", testContract, "-vv")

		// Set the working directory to the subdirectory
		execCmd.Dir = "./levels/"

		// Capture the standard output and standard error of the command
		output, err := execCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", err, output)
		}

		// Print the output of the command to the console
		fmt.Println(string(output))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringP("level", "l", "", "Select a level")
	validateCmd.Flags().StringP("bytecode", "b", "", "The creation bytecode to submit")

	startCmd.MarkFlagRequired("level")
}
