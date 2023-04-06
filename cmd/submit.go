package cmd

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/ethernautdao/evm-runners-cli/internal/config"
	"github.com/spf13/cobra"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit the solution",
	Long:  `Submit the bytecode to the server for processing.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		bytecode, _ := cmd.Flags().GetString("bytecode")
		userID, _ := cmd.Flags().GetString("user_id")
		level, _ := cmd.Flags().GetString("level")

		// load server/auth config
		configStruct, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config")
			return err
		}

		// get level information
		levels, err := config.LoadLevels()
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
		filename := levels[level].FileName
		testContract := levels[level].TestContract

		fmt.Println("Submitting solution for level", level, "with filename", filename)

		// check if bytecode was provided, if not get the bytecode from the huff/sol solution
		if bytecode != "" {
			// check if bytecode is valid
			bytecode = checkValidBytecode(bytecode)

		} else {
			// Check existence of solution files
			_, err1 := os.Stat(fmt.Sprintf("./levels/src/%s.sol", filename))
			_, err2 := os.Stat(fmt.Sprintf("./levels/src/%s.huff", filename))

			if os.IsNotExist(err1) && os.IsNotExist(err2) {
				fmt.Println("No solution file found. Add a solution file or submit bytecode with the --bytecode flag!")
				return nil
			} else if err1 == nil && err2 == nil {
				fmt.Println("More than one solution file found. Delete the one you dont want to submit!")
				return nil
			}

			// .sol solution
			if err1 == nil {
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

				bytecode = checkValidBytecode(bytecodeField["object"].(string))
			}

			// .huff solution
			if err2 == nil {
				// Compile the solution
				execCmd := exec.Command("huffc", fmt.Sprintf("./src/%s.huff", filename), "--bin-runtime")
				execCmd.Dir = "./levels/"
				output, err := execCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("%s: %s", err, output)
				}

				bytecode = checkValidBytecode(string(output))
			}
		}

		fmt.Println("bytecode:", bytecode)

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
			//"level_id": level,
			"level_id": "1", // for now its just id 1. TODO: Change to actual id
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

func checkValidBytecode(bytecode string) string {
	// remove whitespace
	bytecode = strings.TrimSpace(bytecode)
	// remove 0x prefix if present
	bytecode = strings.TrimPrefix(bytecode, "0x")

	// check if bytecode has even length
	if len(bytecode)%2 != 0 {
		fmt.Println("Invalid bytecode length")
		return ""
	}

	// check if bytecode is valid hex
	if _, err := hex.DecodeString(bytecode); err != nil {
		fmt.Println("Invalid bytecode: ", err)
		return ""
	}

	// add 0x prefix again
	bytecode = "0x" + bytecode

	// return sanitized bytecode
	return bytecode
}

func init() {
	rootCmd.AddCommand(submitCmd)

	submitCmd.Flags().StringP("bytecode", "b", "", "The creation bytecode to submit")
	submitCmd.Flags().StringP("user_id", "u", "", "User ID")
	submitCmd.Flags().StringP("level", "l", "", "Level")

	submitCmd.MarkFlagRequired("level")
	submitCmd.MarkFlagRequired("user_id")
}
