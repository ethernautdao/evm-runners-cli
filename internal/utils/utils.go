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
	EVMR_SERVER	string `mapstructure:"EVMR_SERVER"`
	EVMR_TOKEN	string `mapstructure:"EVMR_TOKEN"`
	EVMR_ID		string `mapstructure:"EVMR_ID"`
	EVMR_NAME	string `mapstructure:"EVMR_NAME"`
}

type Level struct {
	ID          string
	File        string
	Name        string
	Test        string
	Description string
}

func LoadConfig() (Config, error) {
	config := Config{EVMR_SERVER: "https://evm-runners.fly.dev/"}

	viper.SetConfigFile(".env")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// load env variables
	config.EVMR_TOKEN = viper.GetString("EVMR_TOKEN")
	config.EVMR_ID = viper.GetString("EVMR_ID")
	config.EVMR_NAME = viper.GetString("EVMR_NAME")

	return config, nil
}

func LoadLevels() (map[string]Level, error) {
	viper.SetConfigFile("./levels/levels.toml")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
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
			ID:          l["id"].(string),
			File:        l["file"].(string),
			Name:        l["name"].(string),
			Test:        l["test"].(string),
			Description: l["description"].(string),
		}
		// Add the new Level struct to the map
		levels[level.Name] = level
	}

	return levels, nil
}

func GetSolves() map[string]string {
	levels, _ := LoadLevels()

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

func CheckValidBytecode(bytecode string) (string, error) {
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

func CheckSolutionFile(File string, langFlag string) (string, error) {
	// Check existence of solution files
	_, err1 := os.Stat(fmt.Sprintf("./levels/src/%s.sol", File))
	_, err2 := os.Stat(fmt.Sprintf("./levels/src/%s.huff", File))

	if os.IsNotExist(err1) && os.IsNotExist(err2) {
		return "", fmt.Errorf("No solution file found. Add a solution file or submit bytecode with the --bytecode flag!")
	} else if err1 == nil && err2 == nil && langFlag == "" {
		return "", fmt.Errorf("More than one solution file found!\nDelete a solution file or use the --lang flag to choose which one to validate.")
	}

	if err1 == nil && (langFlag == "sol" || langFlag == "") {
		return "sol", nil
	} else if err2 == nil && (langFlag == "huff" || langFlag == "") {
		return "huff", nil
	} else {
		return "", fmt.Errorf("Invalid language flag. Please use either 'sol' or 'huff'.")
	}
}
