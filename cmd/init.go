package cmd

import "github.com/ethernautdao/evm-runners-cli/internal/config"

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"reflect"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes EVM Runners",
	Long:  `Initializes EVM Runners by cloning the ethernautdao/evm-runners-levels.git repository into ./levels`,

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Initializing EVM Runners ...")

		// create file
		f, err := os.Create(".env")
		if err != nil {
			fmt.Println("Error creating .env file")
			return err
		}
		// remember to close the file
		defer f.Close()

		// load config
		configStruct, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config")
			return err
		}

		// iterate over struct fields
		v := reflect.ValueOf(configStruct)
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			key := field.Name
			value := v.Field(i).Interface()

			// write key-value pair to file
			_, err := f.WriteString(fmt.Sprintf("%s=%v\n", key, value))
			if err != nil {
				fmt.Println("Error writing to .env file")
				return err
			}
		}

		fmt.Println(".env file created successfully.")

		subdir := "./levels"

		// Check if the subdirectory already exists
		if _, err := os.Stat(subdir); os.IsNotExist(err) {
			// Clone ethernautdao/evm-runners-levels.git repository
			execCmd := exec.Command("git", "clone", "git@github.com:ethernautdao/evm-runners-levels.git", subdir)
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			if err := execCmd.Run(); err != nil {
				fmt.Println("Failed to clone ethernautdao/evm-runners-levels.git")
				return err
			}
			fmt.Println("evm-runners-levels cloned successfully")
		} else {
			fmt.Println("Subdirectory already exists")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
