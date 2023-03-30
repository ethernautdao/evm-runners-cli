package config

import (
	"fmt"
    "github.com/spf13/viper"
)

type Config struct {
    EVMR_PORT   string `mapstructure:"EVMR_PORT"`
    EVMR_SERVER string `mapstructure:"EVMR_SERVER"`
    EVMR_AUTH	string `mapstructure:"EVMR_AUTH"`
}

type Level struct {
    FileName        string
    Contract    	string
	TestContract    string
    Description 	string
}

func LoadConfig() (Config, error) {
    config := Config{EVMR_PORT: "3000", EVMR_SERVER: "http://localhost", EVMR_AUTH: ""}
	
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