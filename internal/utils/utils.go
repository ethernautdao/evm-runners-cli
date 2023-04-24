package utils

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	configFile = ".env"
	levelsFile = "levels.toml"
)

type Config struct {
	EVMR_SERVER     string `mapstructure:"EVMR_SERVER"`
	EVMR_TOKEN      string `mapstructure:"EVMR_TOKEN"`
	EVMR_ID         string `mapstructure:"EVMR_ID"`
	EVMR_NAME       string `mapstructure:"EVMR_NAME"`
	EVMR_LEVELS_DIR string `mapstructure:"EVMR_LEVELS_DIR"`
}

type Level struct {
	ID          string
	File        string
	Name        string
	Test        string
	Description string
}

func LoadConfig() (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		return config, fmt.Errorf("error getting user's home directory: %v", err)
	}

	envFilePath := filepath.Join(usr.HomeDir, ".config", "evm-runners", configFile)
	viper.SetConfigFile(envFilePath)

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading in config file: %v", err)
	}

	// Automatically load environment variables
	viper.AutomaticEnv()

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("error unmarshalling config: %v", err)
	}

	return config, nil
}

func LoadLevels() (map[string]Level, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)

	}

	viper.SetConfigFile(filepath.Join(config.EVMR_LEVELS_DIR, levelsFile))

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading in config file: %v", err)
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
