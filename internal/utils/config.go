package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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
	Contract    string
	Type        string
	Description string
}

func LoadConfig() (Config, error) {
	var config Config

	usr, err := user.Current()
	if err != nil {
		return config, fmt.Errorf("error getting user's home directory: %v", err)
	}

	envFilePath := filepath.Join(usr.HomeDir, ".evm-runners", ".env")
	viper.SetConfigFile(envFilePath)

	// Check if the config file exists before trying to read it
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// print error to run evm-runners init first
		return config, fmt.Errorf("No config file found. Please run 'evm-runners init' first!")
	}

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

func WriteConfig(config Config) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("error getting user's home directory: %v", err)
	}

	envFilePath := filepath.Join(usr.HomeDir, ".evm-runners", ".env")
	viper.SetConfigFile(envFilePath)
	viper.SetConfigType("env")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading in config file: %v", err)
	}

	viper.Set("EVMR_SERVER", config.EVMR_SERVER)
	viper.Set("EVMR_LEVELS_DIR", config.EVMR_LEVELS_DIR)
	viper.Set("EVMR_TOKEN", config.EVMR_TOKEN)
	viper.Set("EVMR_ID", config.EVMR_ID)
	viper.Set("EVMR_NAME", config.EVMR_NAME)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	return nil
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

	// Function to safely convert interface{} to string
	getStringValue := func(val interface{}, fieldName string) string {
		if str, ok := val.(string); ok {
			return str
		}
		fmt.Printf("ERROR: Failed to convert field '%s' to string\nTry running 'evm-runners init' again!\n\n", fieldName)
		return "ERROR"
	}

	// Loop through each level and create a new Level struct
	for _, levelConfig := range levelsConfig {
		l := levelConfig.(map[string]interface{})
		level := Level{
			ID:          getStringValue(l["id"], "id"),
			File:        getStringValue(l["file"], "file"),
			Contract:    getStringValue(l["contract"], "contract"),
			Type:        getStringValue(l["type"], "type"),
			Description: getStringValue(l["description"], "description"),
		}
		// Add the new Level struct to the map
		levels[strings.ToLower(level.Contract)] = level
	}

	return levels, nil
}
