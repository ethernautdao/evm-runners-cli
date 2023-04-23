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
	"path/filepath"
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
	Short: "Submit the solution",
	Long:  `Submit the bytecode to the server for processing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bytecode, _ := cmd.Flags().GetString("bytecode")
		lang, _ := cmd.Flags().GetString("lang")

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return fmt.Errorf("Error loading config: %v", err)
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

		// get filename and test contract of level
		filename := levels[level].File
		testContract := levels[level].Test

		fmt.Println("Submitting solution for level", level, "with filename", filename)

		// check if bytecode was provided, if not get the bytecode from the huff/sol solution
		if bytecode != "" {
			// check if bytecode is valid
			bytecode, err = utils.CheckValidBytecode(bytecode)
			if err != nil {
				return err
			}
		} else {
			solutionType, err := utils.CheckSolutionFile(filename, lang)
			if err != nil {
				return err
			}

			// .sol solution
			if solutionType == "sol" {
				// Compile all contracts
				execCmd := exec.Command("forge", "build")
				execCmd.Dir = config.EVMR_LEVELS_DIR
				output, err := execCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("%s: %s", err, output)
				}

				// Read the JSON file
				file, err := ioutil.ReadFile(filepath.Join(config.EVMR_LEVELS_DIR, "out", fmt.Sprintf("%s.sol", filename), fmt.Sprintf("%s.json", level)))
				if err != nil {
					return fmt.Errorf("error reading JSON file: %v", err)
				}

				// Parse the JSON data
				var data map[string]interface{}
				err = json.Unmarshal([]byte(file), &data)
				if err != nil {
					return fmt.Errorf("error parsing JSON data: %v", err)
				}

				// Extract the "bytecode" field
				bytecodeField := data["bytecode"].(map[string]interface{})

				bytecode, err = utils.CheckValidBytecode(bytecodeField["object"].(string))
				if err != nil {
					return err
				}
			}

			// .huff solution
			if solutionType == "huff" {
				// Compile the solution
				huffPath := filepath.Join("src", fmt.Sprintf("%s.huff", filename))
				execCmd := exec.Command("huffc", huffPath, "--bin-runtime")
				execCmd.Dir = config.EVMR_LEVELS_DIR
				output, err := execCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("%s: %s", err, output)
				}

				bytecode, err = utils.CheckValidBytecode(string(output))
				if err != nil {
					return err
				}
			}
		}

		// Check if solution is correct
		fmt.Println("Validating solution...")

		os.Setenv("BYTECODE", bytecode)
		// Run test
		testContract = testContract + "Base"
		execCmd := exec.Command("forge", "test", "--match-contract", testContract)
		execCmd.Dir = config.EVMR_LEVELS_DIR
		if err = execCmd.Run(); err != nil {
			fmt.Println("Solution is not correct!")
			return nil
		}

		// check if new solution is worse than existing one, if yes ask if user wants to submit anyway

		fmt.Println("Solution is correct! Submitting to the server ...")

		// Create a JSON payload
		payload := map[string]string{
			"bytecode": bytecode,
			"user_id":  config.EVMR_ID,
			"level_id": levels[level].ID,
		}
		jsonPayload, _ := json.Marshal(payload)

		// Make the HTTP request
		url := config.EVMR_SERVER + "submissions"
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.EVMR_TOKEN)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}

		// Parse the response body
		var response SubmitResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		// Print the response
		fmt.Printf("\nSolution for level %s submitted successfully!\nGas: %s, Size: %s\n", level, response.Gas, response.Size)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	submitCmd.Flags().StringP("bytecode", "b", "", "The creation bytecode to submit")
	submitCmd.Flags().StringP("lang", "l", "", "The language of the solution file. Either 'sol' or 'huff'")
}
