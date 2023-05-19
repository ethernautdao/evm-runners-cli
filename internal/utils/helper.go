package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	solutionDir = "src"
)

type SubmissionData struct {
	Gas  string `json:"gas"`
	Size string `json:"size"`
}

// parseOutput function to parse gas and size values
func ParseOutput(output string) (int, int, error) {
	var gasValue int
	var sizeValue int
	var err error
	outputLines := strings.Split(output, "\n")

	for _, line := range outputLines {
		if strings.Contains(line, "_gas") {
			re := regexp.MustCompile(`(μ: )\s*(\d+)`)
			match := re.FindStringSubmatch(line)

			if len(match) > 1 {
				gasValue, err = strconv.Atoi(match[2])
				if err != nil {
					return 0, 0, fmt.Errorf("Error: %v", err)
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
					return 0, 0, fmt.Errorf("Error: %v", err)
				}
			} else {
				fmt.Println("No matching value found")
			}
		}
	}

	return gasValue, sizeValue, nil
}

func GetSolves(levels map[string]Level) map[string]string {
	solves := make(map[string]string)

	config, err := LoadConfig()
	if err != nil {
		return solves

	}

	// Create a custom HTTP client with a 5-second timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for key := range levels {
		url := fmt.Sprintf("%slevels/%s/total", config.EVMR_SERVER, levels[key].ID)
		resp, err := client.Get(url)

		// if the get request errors for some reason, we just set the solve count to 0
		if err != nil {
			solves[levels[key].Contract] = "0"
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			solves[levels[key].Contract] = "0"
			continue
		}

		solves[levels[key].Contract] = string(body)
	}

	return solves
}

func GetBytecodeToValidate(bytecode string, level string, filename string, levelsDir string, lang string) (string, string, error) {
	levels, err := LoadLevels()
	if err != nil {
		return "", "", nil
	}

	// check if bytecode was provided, if not get the bytecode from the huff/sol solution
	if bytecode != "" {
		// check if bytecode is valid
		bytecode, err := validateBytecode(bytecode)
		if err != nil {
			return "", "", err
		}

		return bytecode, "bytecode", nil
	} else {
		solutionType, err := getSolutionType(filename, lang)
		if err != nil {
			return "", "", err
		}

		// .sol solution
		if solutionType == "sol" {
			// Compile all contracts
			execCmd := exec.Command("forge", "build")
			execCmd.Dir = levelsDir
			output, err := execCmd.CombinedOutput()
			if err != nil {
				return "", "", fmt.Errorf("%s: %s", err, output)
			}

			// Read the JSON file
			file, err := ioutil.ReadFile(filepath.Join(levelsDir, "out", fmt.Sprintf("%s.sol", filename), fmt.Sprintf("%s.json", levels[level].Contract)))
			if err != nil {
				return "", "", fmt.Errorf("error reading JSON file: %v", err)
			}

			// Parse the JSON data
			var data map[string]interface{}
			err = json.Unmarshal([]byte(file), &data)
			if err != nil {
				return "", "", fmt.Errorf("error parsing JSON data: %v", err)
			}

			// Extract the "bytecode" field
			bytecodeField := data["bytecode"].(map[string]interface{})

			bytecode, err = validateBytecode(bytecodeField["object"].(string))
			if err != nil {
				return "", "", err
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
				return "", "", fmt.Errorf("%s: %s", err, output)
			}

			bytecode, err = validateBytecode(string(output))
			if err != nil {
				return "", "", err
			}
		}

		// .vy solution
		if solutionType == "vy" {
			// Compile the solution
			vyPath := filepath.Join("src", fmt.Sprintf("%s.vy", filename))
			execCmd := exec.Command("vyper", vyPath)
			execCmd.Dir = levelsDir
			output, err := execCmd.CombinedOutput()
			if err != nil {
				return "", "", fmt.Errorf("%s: %s", err, output)
			}

			bytecode, err = validateBytecode(string(output))
			if err != nil {
				return "", "", err
			}
		}

		return bytecode, solutionType, nil
	}
}

func getSolutionType(file string, langFlag string) (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("error loading config: %v", err)
	}

	// Define the supported languages and their file extensions
	languages := map[string]string{
		"sol":  ".sol",
		"huff": ".huff",
		"vy":   ".vy",
	}

	// Convert the given langFlag to lowercase
	langFlag = strings.ToLower(langFlag)

	// Map additional flags to their corresponding file extensions
	switch langFlag {
	case "solidity":
		langFlag = "sol"
	case "huff":
		langFlag = "huff"
	case "vyper":
		langFlag = "vy"
	}

	// Check if the given langFlag is valid
	if langFlag != "" {
		if _, exists := languages[langFlag]; !exists {
			return "", fmt.Errorf("Invalid language flag. Please use either 'sol', 'huff', or 'vy'.")
		}
	}

	// Check existence of solution files
	var existingFiles []string
	for lang, ext := range languages {
		filePath := filepath.Join(config.EVMR_LEVELS_DIR, solutionDir, file+ext)
		if fileExists(filePath) {
			existingFiles = append(existingFiles, lang)
		}
	}

	// Handle cases with no solution files or multiple solution files
	if len(existingFiles) == 0 {
		return "", fmt.Errorf("No solution file found! Run 'evm-runners start <level>' or submit bytecode with --bytecode")
	} else if langFlag == "" && len(existingFiles) > 1 {
		return "", fmt.Errorf("More than one solution file found!\nDelete a solution file or use --lang to choose which one to validate.")
	}

	// Set langFlag if not provided
	if langFlag == "" {
		langFlag = existingFiles[0]
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

func RunTest(levelsDir string, testContract string, verbose bool) ([]byte, error) {
	// seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random Ethereum address
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		// Handle error
	}
	randAddress := "0x" + hex.EncodeToString(bytes)

	// Generate a random timestamp between 1 Jan 2000 and now
	end := time.Now().Unix()
	randTimestamp := rand.Intn(int(end))

	// Generate a random PrevRandao value
	bytes = make([]byte, 32)
	_, err = rand.Read(bytes)
	if err != nil {
		// Handle error
	}
	randPrevRandao := "0x" + hex.EncodeToString(bytes)

	// initialize the command with common arguments
	execCmd := exec.Command("forge", "test",
		"--block-coinbase", randAddress,
		"--block-timestamp", strconv.Itoa(randTimestamp),
		"--block-number", strconv.Itoa(rand.Intn(17243073)),
		"--block-difficulty", strconv.Itoa(rand.Intn(5875000371)),
		"--block-prevrandao", randPrevRandao,
		"--gas-price", strconv.Itoa(rand.Intn(45014319675)),
		"--base-fee", strconv.Itoa(rand.Intn(45014319675)),
		"--match-contract", testContract)

	// append verbose flag based on verbose variable
	if verbose {
		execCmd.Args = append(execCmd.Args, "-vvvvv")
	} else {
		execCmd.Args = append(execCmd.Args, "-vv")
	}

	execCmd.Dir = levelsDir
	output, err := execCmd.CombinedOutput()

	// Check for errors
	if err != nil {
		return output, err
	}

	return output, nil
}

// fetchSubmissionData function to fetch existing submission data
func FetchSubmissionData(config Config, levelID string) ([]SubmissionData, error) {
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
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error fetching submission data: %s", resp.Status)
	}

	// Parse the response
	var submissions []SubmissionData

	err = json.Unmarshal(body, &submissions)
	if err != nil {
		return nil, fmt.Errorf("Error parsing the response: %v", err)
	}

	return submissions, nil
}
