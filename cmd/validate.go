package cmd

import (
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate <level>",
	Short: "Validates a level",
	Long: `Validates a level by running the predefined Foundry tests against 
the solution file or against the provided bytecode, if the bytecode -b flag is set.

The resulting codesize score is determined by the result of 'test_<level_id>_size',
and the gas score is determined by the µ value of the 'test_<level_id>_gas' fuzz test.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		bytecode, _ := cmd.Flags().GetString("bytecode")
		lang, _ := cmd.Flags().GetString("lang")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if len(args) == 0 {
			return fmt.Errorf("Please provide a level\n")
		}
		level := strings.ToLower(args[0])

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

		bytecode, solutionType, err := utils.GetBytecodeToValidate(bytecode, level, filename, config.EVMR_LEVELS_DIR, lang)
		if err != nil {
			return err
		}

		// Check if solution is correct
		fmt.Printf("Validating solution...\n\n")

		os.Setenv("BYTECODE", bytecode)

		// Run test
		testContract := levels[level].Name + "TestBase"

		// run forge test based on verbose flag
		var execCmd *exec.Cmd
		if verbose {
			var userTestContract string
			switch solutionType {
			case "sol":
				userTestContract = levels[level].Name + "TestSol"
			case "huff":
				userTestContract = levels[level].Name + "TestHuff"
			case "vyper":
				userTestContract = levels[level].Name + "TestVyper"
			case "bytecode":
				userTestContract = testContract
			}

			// show user which command is run
			fmt.Printf("To test the solution yourself, run 'forge test --mc %s -vvvvv' in %s\n\n", userTestContract, config.EVMR_LEVELS_DIR)
			execCmd = exec.Command("forge", "test", "--match-contract", testContract, "-vvvvv")
		} else {
			execCmd = exec.Command("forge", "test", "--match-contract", testContract, "-vv")
		}

		execCmd.Dir = config.EVMR_LEVELS_DIR
		output, err := execCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s", output)
		}

		// Print the output of the command to the console
		fmt.Println(string(output))

		// Parse the output to get gas and size values
		gasValue, sizeValue, err := utils.ParseOutput(string(output))
		if err != nil {
			return err
		}

		// Print the gas and size values
		fmt.Printf("Solution is correct! Gas: %d, Size: %d\n", gasValue, sizeValue)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringP("bytecode", "b", "", "The creation bytecode to submit")
	validateCmd.Flags().StringP("lang", "l", "", "The language of the solution file. Either 'sol' or 'huff'")
	validateCmd.Flags().BoolP("verbose", "v", false, "Verbose output, shows stack and setup traces")
}
