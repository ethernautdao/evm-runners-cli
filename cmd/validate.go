package cmd

import (
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		// get level information
		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("Error loading levels: %v", err)
		}

		// check if level exists
		if _, ok := levels[level]; !ok {
			return fmt.Errorf("Invalid level: %v", level)
		}

		// get filename and test contract of level
		filename := levels[level].File

		bytecode, err = utils.GetBytecodeToValidate(bytecode, level, filename, config.EVMR_LEVELS_DIR, lang)
		if err != nil {
			return err
		}

		// Check if solution is correct
		fmt.Printf("Validating solution...\n\n")

		os.Setenv("BYTECODE", bytecode)

		// Run test
		testContract := level + "TestBase"
		execCmd := exec.Command("forge", "test", "--match-contract", testContract)
		execCmd.Dir = config.EVMR_LEVELS_DIR
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
	validateCmd.Flags().StringP("lang", "l", "", "The language of the solution file. Either 'sol' or 'huff'")
}
