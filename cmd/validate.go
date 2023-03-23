package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
    Use:   "validate",
    Short: "Validates a level",
    Long:  `Validates a level by running the predefined Foundry tests against the solution.`,

    RunE: func(cmd *cobra.Command, args []string) error {
    	level, _ := cmd.Flags().GetString("level")

		testContract := getTestContract(level);

		// Check if the level is valid
		if testContract == "" {
			fmt.Println("Invalid level")
			return nil
		}

		// Set the levels path
		subDir := "./levels/"

		// Create the command to be run in the subdirectory
		execCmd := exec.Command("forge", "test", "--match-contract", level, "-vv")

		// Set the working directory to the subdirectory
		execCmd.Dir = subDir

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

func getTestContract(level string) string {
	switch level {
	case "S01E01-Average":
		return "AverageTest"
	case "Average":
		return "AverageTest"
	case "S01E01":
		return "AverageTest"
	}

	return ""
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("level", "l", "", "Select a level")
	startCmd.MarkFlagRequired("level")
}