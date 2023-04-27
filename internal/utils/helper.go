package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	solutionDir = "src"
)

func GetSolves() map[string]string {
	levels, err := LoadLevels()
	if err != nil {
		return nil
	}

	solves := make(map[string]string)

	for key := range levels {
		url := fmt.Sprintf("https://evm-runners.fly.dev/levels/%s/total", levels[key].ID)
		resp, err := http.Get(url)

		// if the get request errors for some reason, we just set the solve count to 0
		if err != nil {
			//fmt.Printf("Error fetching submission count for level %s: %v\n", levels[key].Name, err)
			solves[levels[key].Name] = "0"
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//fmt.Printf("Error reading response body for level %s: %v\n",  levels[key].Name, err)
			solves[levels[key].Name] = "0"
			continue
		}

		solves[levels[key].Name] = string(body)
	}

	return solves
}

func GetBytecodeToValidate(bytecode string, level string, filename string, levelsDir string, lang string) (string, error) {

	// check if bytecode was provided, if not get the bytecode from the huff/sol solution
	if bytecode != "" {
		// check if bytecode is valid
		bytecode, err := validateBytecode(bytecode)
		if err != nil {
			return "", err
		}

		return bytecode, nil
	} else {
		solutionType, err := getSolutionFile(filename, lang)
		if err != nil {
			return "", err
		}

		// .sol solution
		if solutionType == "sol" {
			// Compile all contracts
			execCmd := exec.Command("forge", "build")
			execCmd.Dir = levelsDir
			output, err := execCmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("%s: %s", err, output)
			}

			// Read the JSON file
			file, err := ioutil.ReadFile(filepath.Join(levelsDir, "out", fmt.Sprintf("%s.sol", filename), fmt.Sprintf("%s.json", level)))
			if err != nil {
				return "", fmt.Errorf("error reading JSON file: %v", err)
			}

			// Parse the JSON data
			var data map[string]interface{}
			err = json.Unmarshal([]byte(file), &data)
			if err != nil {
				return "", fmt.Errorf("error parsing JSON data: %v", err)
			}

			// Extract the "bytecode" field
			bytecodeField := data["bytecode"].(map[string]interface{})

			bytecode, err = validateBytecode(bytecodeField["object"].(string))
			if err != nil {
				return "", err
			}
		}

		// .huff solution
		if solutionType == "huff" {
			// Compile the solution
			huffPath := filepath.Join("src", fmt.Sprintf("%s.huff", filename))
			execCmd := exec.Command("huffc", huffPath, "--bytecode")
			execCmd.Dir = levelsDir
			output, err := execCmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("%s: %s", err, output)
			}

			bytecode, err = validateBytecode(string(output))
			if err != nil {
				return "", err
			}
		}
		return bytecode, nil
	}
}

func getSolutionFile(file string, langFlag string) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("error loading config: %v", err)
	}

	// Check existence of solution files
	solFile := filepath.Join(config.EVMR_LEVELS_DIR, solutionDir, file+".sol")
	huffFile := filepath.Join(config.EVMR_LEVELS_DIR, solutionDir, file+".huff")

	solExists := fileExists(solFile)
	huffExists := fileExists(huffFile)

	if !solExists && !huffExists {
		return "", fmt.Errorf("No solution file found! Run 'evm-runners start <level>' or submit bytecode with --bytecode")
	}

	if langFlag != "" && langFlag != "sol" && langFlag != "huff" {
		return "", fmt.Errorf("Invalid language flag. Please use either 'sol' or 'huff'.")
	}

	if langFlag == "" {
		if solExists && huffExists {
			return "", fmt.Errorf("More than one solution file found!\nDelete a solution file or use --lang to choose which one to validate.")
		}
		if solExists {
			langFlag = "sol"
		} else {
			langFlag = "huff"
		}
	}

	if (langFlag == "sol" && !solExists) || (langFlag == "huff" && !huffExists) {
		return "", fmt.Errorf("Solution file not found for the specified language flag.")
	}

	return langFlag, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func validateBytecode(bytecode string) (string, error) {
	// remove whitespace
	bytecode = strings.TrimSpace(bytecode)
	// remove 0x prefix if present
	bytecode = strings.TrimPrefix(bytecode, "0x")

	// check if bytecode has even length
	if len(bytecode)%2 != 0 {
		return "", fmt.Errorf("Invalid bytecode length")
	}

	// check if bytecode is valid hex
	if _, err := hex.DecodeString(bytecode); err != nil {
		return "", fmt.Errorf("Invalid bytecode: %v", err)

	}

	// add 0x prefix again
	bytecode = "0x" + bytecode

	// return sanitized bytecode
	return bytecode, nil
}
