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
	Short: "Submit a solution for a level to the server",
	Long: `Submit a solution for a level to the server.

This command performs the following steps:

1. Compiling the solution
2. Validating the solution's bytecode 
3. If successful, sending the bytecode to the server

The displayed scores can differ slightly from the final scores on the leaderboard.`,

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
			return fmt.Errorf("Please authorize first with 'evm-runners auth discord'\n")
		}

		if len(args) == 0 {
			return fmt.Errorf("Please provide a level\n")
		}
		level := strings.ToLower(args[0])

		// get level information
		levels, err := utils.LoadLevels()
		if err != nil {
			return fmt.Errorf("error loading levels: %v", err)
		}

		// check if level exists
		if _, ok := levels[level]; !ok {
			return fmt.Errorf("Invalid level: %v\n", level)
		}

		// get filename of level
		filename := levels[level].File

		bytecode, solutionType, err := utils.GetBytecodeToValidate(bytecode, level, filename, config.EVMR_LEVELS_DIR, lang)
		if err != nil {
			return err
		}

		// Check if solution is correct
		fmt.Printf("Validating solution for level '%s' ...\n", level)

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
		submissions, err := utils.FetchSubmissionData(config)
		if err != nil {
			return err
		}

		if len(submissions) > 0 {
			// Compare new solution's gas and size with existing submission
			var existingGas int
			var existingSize int

			for _, item := range submissions {
				if level == strings.ToLower(item.LevelName) {
					existingGas, _ = strconv.Atoi(item.Gas)
					existingSize, _ = strconv.Atoi(item.Size)
				}
			}

			// If gas and size score is worse than existing one, skip submission
			if gasValue >= existingGas && sizeValue >= existingSize {
				fmt.Printf("\nWarning: Submission skipped!\nThe gas and size scores are either worse or equal to the existing ones (gas: %d, size: %d).\n", existingGas, existingSize)
				return nil
			}
		}

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
			return fmt.Errorf("error sending the request: %v", err)
		}
		defer resp.Body.Close()

		// Check for errors in the response
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("http request failed with status: %s", resp.Status)
		}

		// Decode the JSON response as an array of objects
		var response []map[string]interface{}
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&response); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		// Extract the gas and size rank from the first object in the array
		gasRank, _ := response[0]["gas_rank"].(string)
		sizeRank, _ := response[0]["size_rank"].(string)

		fmt.Printf("\nSolution for level '%s' submitted successfully!\n\n", level)
		fmt.Printf("Size leaderboard: #%s\nGas leaderboard: #%s\n", sizeRank, gasRank)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Flags
	submitCmd.Flags().StringP("bytecode", "b", "", "The bytecode of the solution")
	submitCmd.Flags().StringP("lang", "l", "", "The language of the solution file (sol, huff, vyper)")
}
