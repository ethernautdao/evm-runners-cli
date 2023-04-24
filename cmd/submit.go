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
	"regexp"
	"strconv"
	"strings"
)

type SubmitResponse struct {
	ID       string `json:"id"`
	LevelID  int    `json:"level_id"`
	UserID   int    `json:"user_id"`
	Bytecode string `json:"bytecode"`
	Gas      string `json:"gas"`
	Size     string `json:"size"`
}

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit <level>",
	Short: "Submits a solution for a level to the server",
	Long: `Submits a solution for a level to the server by 
	
1. Validating if the solution is correct
2. Checking if an existing solution will be overwritten
3. Sending the bytecode to the server`,

	RunE: func(cmd *cobra.Command, args []string) error {
		bytecode, _ := cmd.Flags().GetString("bytecode")
		lang, _ := cmd.Flags().GetString("lang")

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}

		// check if user authenticated
		if config.EVMR_TOKEN == "" {
			return fmt.Errorf("Please authorize first with 'evm-runners auth discord'")
		}

		if len(args) == 0 {
			return fmt.Errorf("Please provide a level\n")
		}
		level := args[0]

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

		fmt.Println("Submitting solution for level", level, "...")

		bytecode, err = utils.GetBytecodeToValidate(bytecode, level, filename, config.EVMR_LEVELS_DIR, lang)
		if err != nil {
			return err
		}

		// Check if solution is correct
		fmt.Println("Validating solution...")

		os.Setenv("BYTECODE", bytecode)

		// Run test
		testContract := levels[level].Test + "Base"
		execCmd := exec.Command("forge", "test", "--match-contract", testContract)
		execCmd.Dir = config.EVMR_LEVELS_DIR
		output, err := execCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Solution is not correct!")
			return nil
		}

		// Parse the output to get gas and size values
		var gasValue int
		var sizeValue int
		outputStr := string(output)
		outputLines := strings.Split(outputStr, "\n")
		for _, line := range outputLines {
			if strings.Contains(line, "_gas") {
				re := regexp.MustCompile(`(μ|~:)\s*(\d+)`)
				match := re.FindStringSubmatch(line)

				if len(match) > 0 {
					gasValue, err = strconv.Atoi(match[2])
					if err != nil {
						fmt.Printf("Error: %s\n", err.Error())
					}
				} else {
					fmt.Println("No matching value found")
				}
			}
			if strings.Contains(line, "Contract size:") {

				re := regexp.MustCompile(`Contract size:\s*(\d+)`)
				match := re.FindStringSubmatch(line)

				if len(match) > 1 {
					sizeValue, err = strconv.Atoi(match[1])
					if err != nil {
						fmt.Printf("Error: %s\n", err.Error())
					}
				} else {
					fmt.Println("No matching value found")
				}
			}
		}

		fmt.Printf("Solution is correct! Gas: %d, Size: %d\nSubmitting to the server...\n", gasValue, sizeValue)

		// Fetch existing submission data
		url := config.EVMR_SERVER + "submissions/user/" + levels[level].ID
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+config.EVMR_TOKEN)

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

		// Parse the response
		type SubmissionData struct {
			Gas  string `json:"gas"`
			Size string `json:"size"`
		}
		var submissions []SubmissionData

		err = json.Unmarshal(body, &submissions)
		if err != nil {
			return fmt.Errorf("Error parsing the response: %v", err)
		}

		// Compare new solution's gas and size with existing submission
		var existingGas int
		var existingSize int
		if len(submissions) > 0 {
			floatGas, _ := strconv.ParseFloat(submissions[0].Gas, 64)
			existingGas = int(floatGas)
			existingSize, _ = strconv.Atoi(submissions[0].Size)

			if gasValue >= existingGas || sizeValue >= existingSize {
				fmt.Printf("Warning: Gas (%d) or size (%d) of the new solution is higher or equal to the existing submission (gas: %d, size: %d).\n", gasValue, sizeValue, existingGas, existingSize)
				fmt.Print("Do you want to submit anyway? (y/n): ")
				var overwrite string
				fmt.Scanln(&overwrite)
				if overwrite != "y" && overwrite != "Y" {
					fmt.Println("Submission cancelled.")
					return nil
				}
			}
		}

		// Create a JSON payload
		payload := map[string]string{
			"bytecode": bytecode,
			"user_id":  config.EVMR_ID,
			"level_id": levels[level].ID,
		}
		jsonPayload, _ := json.Marshal(payload)

		// Make the HTTP request
		url = config.EVMR_SERVER + "submissions"
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.EVMR_TOKEN)

		// Send the request
		client = &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("Error sending the request: %v", err)
		}
		defer resp.Body.Close()

		// Read the response
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error reading the response: %v", err)
		}

		// Check for errors in the response
		if resp.StatusCode != 200 {
			return fmt.Errorf("Error submitting solution: %s", body)
		}

		// Parse the response
		var res SubmitResponse
		err = json.Unmarshal(body, &res)
		if err != nil {
			return fmt.Errorf("Error parsing the response: %v", err)
		}

		// Print the response
		fmt.Printf("\nSolution for level %s submitted successfully!", level)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Flags
	submitCmd.Flags().StringP("bytecode", "b", "", "The bytecode of the solution")
	submitCmd.Flags().StringP("lang", "l", "", "The programming language of the solution (e.g. 'huff' or 'sol')")
}
