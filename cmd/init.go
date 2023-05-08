package cmd

import (
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes EVM Runners",
	Long: `Initializes EVM Runners by

1. Cloning the ethernautdao/evm-runners-levels.git repository into ./evm-runners-levels
2. Creating a .env file in ~/.config/evm-runners/`,

	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Initializing EVM Runners ...")

		// Get user's home directory
		usr, err := user.Current()
		if err != nil {
			return fmt.Errorf("error getting user's home directory: %v", err)
		}

		// Get absolute path for evm-runners-levels
		subdir, err := filepath.Abs("evm-runners-levels")
		if err != nil {
			return fmt.Errorf("error getting absolute path for evm-runners-levels: %v", err)
		}

		fmt.Printf("\nCloning ethernautdao/evm-runners-levels.git ...\n")

		if _, err := os.Stat(subdir); os.IsNotExist(err) {
			execCmd := exec.Command("git", "clone", "https://github.com/ethernautdao/evm-runners-levels.git", subdir)
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			if err := execCmd.Run(); err != nil {
				return fmt.Errorf("error cloning ethernautdao/evm-runners-levels.git: %v", err)
			}
			fmt.Println("evm-runners-levels cloned successfully")
		} else {
			fmt.Println("evm-runners-levels already exists")
		}

		envDirPath := filepath.Join(usr.HomeDir, ".evm-runners")
		envFilePath := filepath.Join(envDirPath, ".env")

		fmt.Printf("\nCreating .env file at %s ...\n", envFilePath)

		// Set the fields in the config struct
		config := utils.Config{
			EVMR_SERVER:     "https://evm-runners.fly.dev/",
			EVMR_LEVELS_DIR: subdir,
		}

		// Write or overwrite .env file
		createEnvFile := func() error {
			f, err := os.Create(envFilePath)
			if err != nil {
				return fmt.Errorf("error creating .env file: %v", err)
			}
			defer f.Close()

			err = utils.WriteConfig(config)
			if err != nil {
				return fmt.Errorf("error writing to .env file: %v", err)
			}
			return nil
		}

		// Check if the .env file exists
		if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
			// Create the directory if it doesn't exist
			if _, err := os.Stat(envDirPath); os.IsNotExist(err) {
				if err := os.MkdirAll(envDirPath, 0755); err != nil {
					return fmt.Errorf("error creating directory %s: %v", envDirPath, err)
				}
			}

			err = createEnvFile()
			if err != nil {
				return err
			}
			fmt.Println(".env file created successfully.")
		} else {
			fmt.Printf(".env file already exists at the destination.\nDo you want to update it? (y/n): ")
			var overwrite string
			fmt.Scanln(&overwrite)
			if overwrite != "y" && overwrite != "Y" {
				fmt.Println("\nNot updating .env file")
			} else {
				fmt.Printf("\nUpdating .env file at %s ...\n", envFilePath)

				// Load existing config
				config, err := utils.LoadConfig()
				if err != nil {
					return err
				}

				// replace subdir in config
				config.EVMR_LEVELS_DIR = subdir

				// Update config
				err = utils.WriteConfig(config)
				if err != nil {
					return fmt.Errorf("error writing to .env file: %v", err)
				}
				fmt.Println(".env file updated successfully.")
			}
		}

		fmt.Println("\nevm-runners initialized successfully!\nSee 'evm-runners --help' for a list of all available commands.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
