package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit <level>",
	Short: "Submits a solution for a level to the server",
	Long: `Submits a solution for a level to the server by 

1. Compiling the solution
2. Validating the solution's bytecode 
3. Sending the bytecode to the server

The resulting codesize score is determined by the result
of 'test_<level_id>_size', and the gas score is determined 
by the µ value of the 'test_<level_id>_gas' fuzz test.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		bytecode, _ := cmd.Flags().GetString("bytecode")
		lang, _ := cmd.Flags().GetString("lang")

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		// check if user authenticated
		if config.EVMR_TOKEN == "" {
			return fmt.Errorf("Please authorize first with 'evm-runners auth discord'")
		}

		if len(args) == 0 {
			return fmt.Errorf("Please provide a level\n")
		}
		level := strings.ToLower(args[0])

		// get level information
		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("Error loading levels: %v", err)
		}

		// check if level exists
		if _, ok := levels[level]; !ok {
			return fmt.Errorf("Invalid level: %v", level)
		}

		// get filename of level
		filename := levels[level].File

		bytecode, solutionType, err := utils.GetBytecodeToValidate(bytecode, level, filename, config.EVMR_LEVELS_DIR, lang)
		if err != nil {
			return err
		}

		// Check if solution is correct
		fmt.Printf("\nValidating solution for level '%s' ...\n", level)

		os.Setenv("BYTECODE", bytecode)

		// Run test
		testContract := levels[level].Contract + "TestBase"
		output, err := utils.RunTest(config.EVMR_LEVELS_DIR, testContract, false)
		if err != nil {
			fmt.Println("Solution is not correct!")
			return nil
		}

		// Parse the output to get gas and size values
		gasValue, sizeValue, err := utils.ParseOutput(string(output))
		if err != nil {
			return err
		}

		fmt.Printf("Solution is correct! Gas: %d, Size: %d\nNote: The final score can be slightly different.\n", gasValue, sizeValue)

		// Fetch existing submission data
		submissions, err := utils.FetchSubmissionData(config, levels[level].ID)
		if err != nil {
			return err
		}

		// Compare new solution's gas and size with existing submission
		var existingGas int
		var existingSize int
		if len(submissions) > 0 {
			floatGas, _ := strconv.ParseFloat(submissions[0].Gas, 64)
			existingGas = int(floatGas)
			existingSize, _ = strconv.Atoi(submissions[0].Size)

			if gasValue >= existingGas && sizeValue >= existingSize {
				fmt.Printf("\nWarning: Solution skipped!\nGas and size score is larger than existing one (gas: %d, size: %d).\n", existingGas, existingSize)
				return nil
			}
		}

		// If solution is better than existing one, submit it

		// Create a JSON payload
		payload := map[string]string{
			"bytecode": bytecode,
			"type":     solutionType,
			"user_id":  config.EVMR_ID,
			"level_id": levels[level].ID,
		}
		jsonPayload, _ := json.Marshal(payload)

		// Make the HTTP request
		url := config.EVMR_SERVER + "submissions"
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.EVMR_TOKEN)

		// Send the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("Error sending the request: %v", err)
		}
		defer resp.Body.Close()

		// Check for errors in the response
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Error submitting solution: %s", resp.Status)
		}

		// Decode the JSON response as an array of objects
		var response []map[string]interface{}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&response); err != nil {
			return fmt.Errorf("Error decoding response: %v", err)
		}

		// Extract the gas and size rank from the first object in the array
		gasRank, _ := response[0]["gas_rank"].(string)
		sizeRank, _ := response[0]["size_rank"].(string)

		fmt.Printf("\nSolution for level '%s' submitted successfully!\n\n", level)
		fmt.Printf("Size leaderboard: #%s\nGas leaderboard: #%s\n\n", sizeRank, gasRank)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Flags
	submitCmd.Flags().StringP("bytecode", "b", "", "The bytecode of the solution")
	submitCmd.Flags().StringP("lang", "l", "", "The programming language of the solution (e.g. 'huff' or 'sol')")
}
