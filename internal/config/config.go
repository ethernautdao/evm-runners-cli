package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
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

func GetSolves() (map[string]string, error) {
	levels, err := LoadLevels()
	if err != nil {
		fmt.Println("Error loading levels")
		return nil, err
	}

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

	return solves, nil
}
