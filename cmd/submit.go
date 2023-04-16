package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit <level>",
	Short: "Submit the solution",
	Long:  `Submit the bytecode to the server for processing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bytecode, _ := cmd.Flags().GetString("bytecode")
		userID, _ := cmd.Flags().GetString("user_id")
		lang, _ := cmd.Flags().GetString("lang")

		// load server/auth config
		configStruct, err := utils.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config")
			return err
		}

		if len(args) == 0 {
			return fmt.Errorf("Please provide a level\n")
		}
		level := args[0]

		// get level information
		levels, err := utils.LoadLevels()
		if err != nil {
			fmt.Println("Error loading levels")
			return err
		}

		// check if level exists
		if _, ok := levels[level]; !ok {
			fmt.Println("Invalid level")
			return nil
		}

		// get filename and test contract of level
		filename := levels[level].File
		testContract := levels[level].Test

		fmt.Println("Submitting solution for level", level, "with filename", filename)

		// check if bytecode was provided, if not get the bytecode from the huff/sol solution
		if bytecode != "" {
			// check if bytecode is valid
			bytecode = utils.CheckValidBytecode(bytecode)

		} else {
			solutionType := utils.CheckSolutionFile(filename, lang)
			if solutionType == "nil" {
				return nil
			}

			// .sol solution
			if solutionType == "sol" {
				// Compile all contracts
				execCmd := exec.Command("forge", "build")
				execCmd.Dir = "./levels/"
				output, err := execCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("%s: %s", err, output)
				}

				// Read the JSON file
				file, err := ioutil.ReadFile(fmt.Sprintf("./levels/out/%s.sol/%s.json", filename, level))
				if err != nil {
					panic(err)
				}

				// Parse the JSON data
				var data map[string]interface{}
				err = json.Unmarshal([]byte(file), &data)
				if err != nil {
					panic(err)
				}

				// Extract the "bytecode" field
				bytecodeField := data["bytecode"].(map[string]interface{})

				bytecode = utils.CheckValidBytecode(bytecodeField["object"].(string))
			}

			// .huff solution
			if solutionType == "huff" {
				// Compile the solution
				execCmd := exec.Command("huffc", fmt.Sprintf("./src/%s.huff", filename), "--bin-runtime")
				execCmd.Dir = "./levels/"
				output, err := execCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("%s: %s", err, output)
				}

				bytecode = utils.CheckValidBytecode(string(output))
			}
		}

		//fmt.Println("bytecode:", bytecode)

		// Check if solution is correct
		fmt.Println("Validating solution...")

		os.Setenv("BYTECODE", bytecode)
		// Run test
		execCmd := exec.Command("forge", "test", "--match-contract", testContract)
		execCmd.Dir = "./levels/"
		if err = execCmd.Run(); err != nil {
			fmt.Println("Solution is not correct!")
			return nil
		}

		fmt.Println("Solution is correct! Submitting to server ...")

		// Create a JSON payload
		payload := map[string]string{
			"bytecode": bytecode,
			"user_id":  userID,
			"level_id": levels[level].ID,
		}
		jsonPayload, _ := json.Marshal(payload)

		// Make the HTTP request
		url := configStruct.EVMR_SERVER + "submissions"
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// Print the response
		fmt.Println("Response Status:", resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("Response Body:", string(body))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	submitCmd.Flags().StringP("bytecode", "b", "", "The creation bytecode to submit")
	submitCmd.Flags().StringP("user_id", "u", "", "User ID")
	submitCmd.Flags().StringP("lang", "l", "", "The language of the solution file. Either 'sol' or 'huff'")

	submitCmd.MarkFlagRequired("user_id")
}
