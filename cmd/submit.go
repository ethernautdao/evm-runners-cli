package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SubmissionData struct {
	Gas  string `json:"gas"`
	Size string `json:"size"`
}

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
		fmt.Println("\nValidating solution for level", level, "...")

		os.Setenv("BYTECODE", bytecode)

		// Run test
		testContract := levels[level].Name + "TestBase"
		execCmd := exec.Command("forge", "test", "--match-contract", testContract)
		execCmd.Dir = config.EVMR_LEVELS_DIR
		output, err := execCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Solution is not correct!")
			return nil
		}

		// Parse the output to get gas and size values
		gasValue, sizeValue, err := utils.ParseOutput(string(output))
		if err != nil {
			return err
		}

		fmt.Printf("Solution is correct! Gas: %d, Size: %d\n\nSubmitting to the server...\n", gasValue, sizeValue)

		// Fetch existing submission data
		submissions, err := fetchSubmissionData(config, levels[level].ID)
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

			if gasValue >= existingGas || sizeValue >= existingSize {
				fmt.Printf("\nWarning: Gas (%d) or size (%d) of the new solution is higher or equal to the existing solution (gas: %d, size: %d).\n", gasValue, sizeValue, existingGas, existingSize)
				fmt.Println("Only the better score will be replaced.")
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
			return fmt.Errorf("Error sending the request: %v", err)
		}
		defer resp.Body.Close()

		// Read the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading the response: %v", err)
		}

		// Check for errors in the response
		if resp.StatusCode != 200 {
			return fmt.Errorf("Error submitting solution: %s", body)
		}

		fmt.Printf("\nSolution for level '%s' submitted successfully!\n", level)

		return nil
	},
}

// fetchSubmissionData function to fetch existing submission data
func fetchSubmissionData(config utils.Config, levelID string) ([]SubmissionData, error) {
	url := config.EVMR_SERVER + "submissions/user/" + levelID
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+config.EVMR_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending the request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading the response: %v", err)
	}

	// Check for errors in the response
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error submitting solution (status code %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var submissions []SubmissionData

	err = json.Unmarshal(body, &submissions)
	if err != nil {
		return nil, fmt.Errorf("Error parsing the response: %v", err)
	}

	return submissions, nil
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Flags
	submitCmd.Flags().StringP("bytecode", "b", "", "The bytecode of the solution")
	submitCmd.Flags().StringP("lang", "l", "", "The programming language of the solution (e.g. 'huff' or 'sol')")
}
