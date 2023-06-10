package cmd

import (
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate <level>",
	Short: "Validates a level",
	Long: `Validates a level by running the predefined Foundry tests 
against the solution file or against the provided bytecode,
if the bytecode -b flag is set.

The resulting codesize score is determined by the result
of 'test_<level_id>_size', and the gas score is determined 
by the µ value of the 'test_<level_id>_gas' fuzz test.`,

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
		testContract := levels[level].Contract + "TestBase"

		output, err := utils.RunTest(config.EVMR_LEVELS_DIR, testContract, verbose)
		if err != nil {
			// print the output of forge test
			fmt.Printf("%s", output)

			// if verbose == true, show the test command to the user
			if verbose {
				var userTestContract string
				switch solutionType {
				case "sol":
					userTestContract = levels[level].Contract + "TestSol"
				case "huff":
					userTestContract = levels[level].Contract + "TestHuff"
				case "vyper":
					userTestContract = levels[level].Contract + "TestVyper"
				case "bytecode":
					userTestContract = testContract
				}

				fmt.Printf("\nTo test the solution yourself, run 'forge test --mc %s -vvvv' in %s\n\n", userTestContract, config.EVMR_LEVELS_DIR)
			}

			return nil
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
		fmt.Printf("To submit it, run 'evm-runners submit %s'\n", level)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringP("bytecode", "b", "", "The creation bytecode to submit")
	validateCmd.Flags().StringP("lang", "l", "", "The language of the solution file. Either 'sol' or 'huff'")
	validateCmd.Flags().BoolP("verbose", "v", false, "Verbose output, shows stack traces of all tests")
}
