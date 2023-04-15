package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Config struct {
	EVMR_SERVER string `mapstructure:"EVMR_SERVER"`
	EVMR_AUTH   string `mapstructure:"EVMR_AUTH"`
}

type Level struct {
	ID           string
	FileName     string
	Contract     string
	TestContract string
	Description  string
}

func LoadConfig() (Config, error) {
	config := Config{EVMR_SERVER: "https://evm-runners.fly.dev/", EVMR_AUTH: ""}

	viper.SetConfigFile(".env")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return config, err
	}

	config.EVMR_AUTH = viper.GetString("EVMR_AUTH")

	return config, nil
}

func LoadLevels() (map[string]Level, error) {
	viper.SetConfigFile("./levels/levels.toml")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
		return nil, err
	}

	// Get the levels array from the config file
	levelsConfig := viper.Get("levels").([]interface{})

	// Create a map of Level structs
	levels := make(map[string]Level)

	// Loop through each level and create a new Level struct
	for _, levelConfig := range levelsConfig {
		l := levelConfig.(map[string]interface{})
		level := Level{
			ID:           l["id"].(string),
			FileName:     l["filename"].(string),
			Contract:     l["contract"].(string),
			TestContract: l["testcontract"].(string),
			Description:  l["description"].(string),
		}
		// Add the new Level struct to the map
		levels[level.Contract] = level
	}

	return levels, nil
}

func GetSolves() map[string]string {
	levels, _ := LoadLevels()

	solves := make(map[string]string)

	for key, _ := range levels {
		url := fmt.Sprintf("https://evm-runners.fly.dev/levels/%s/total", levels[key].ID)
		resp, err := http.Get(url)

		// if the get request errors for some reason, we just set the solve count to 0
		if err != nil {
			//fmt.Printf("Error fetching submission count for level %s: %v\n", levels[key].Contract, err)
			solves[levels[key].Contract] = "0"
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			//fmt.Printf("Error reading response body for level %s: %v\n",  levels[key].Contract, err)
			solves[levels[key].Contract] = "0"
			continue
		}

		solves[levels[key].Contract] = string(body)
	}

	return solves
}

func CheckValidBytecode(bytecode string) string {
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

// TODO: return costum error?
func CheckSolutionFile(filename string, langFlag string) string {
	// Check existence of solution files
	_, err1 := os.Stat(fmt.Sprintf("./levels/src/%s.sol", filename))
	_, err2 := os.Stat(fmt.Sprintf("./levels/src/%s.huff", filename))

	if os.IsNotExist(err1) && os.IsNotExist(err2) {
		fmt.Println("No solution file found. Add a solution file or submit bytecode with the --bytecode flag!")
		return "nil"
	} else if err1 == nil && err2 == nil && langFlag == "" {
		fmt.Println("More than one solution file found!\nDelete a solution file or use the --lang flag to choose which one to validate.")
		return "nil"
	}

	if err1 == nil && (langFlag == "sol" || langFlag == "") {
		return "sol"
	} else if err2 == nil && (langFlag == "huff" || langFlag == "") {
		return "huff"
	} else {
		fmt.Println("Invalid language flag. Please use either 'sol' or 'huff'.")
		return "nil"
	}
}
