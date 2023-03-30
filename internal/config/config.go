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

func LoadLevels() ([]Level, error) {
	viper.SetConfigFile("./levels/levels.toml")

	// Read the config file
    if err := viper.ReadInConfig(); err != nil {
        fmt.Println("Error reading config file:", err)
        return nil, err
    }

    // Get the levels array from the config file
    levelsConfig := viper.Get("levels").([]interface{})
    
    // Create a slice of Level structs
    levels := make([]Level, 0, len(levelsConfig))
    
    // Loop through each level and create a new Level struct
    for _, levelConfig := range levelsConfig {
        l := levelConfig.(map[string]interface{})
        level := Level{
            FileName:     l["filename"].(string),
            Contract:     l["contract"].(string),
            TestContract: l["testcontract"].(string),
            Description:  l["description"].(string),
        }
        // Append the new Level struct to the slice
        levels = append(levels, level)
    }

	return levels, nil
}