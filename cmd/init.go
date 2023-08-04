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
	Short: "Initialize evm-runners",
	Long: `Initialize evm-runners by

1. Cloning the ethernautdao/evm-runners-levels.git repository into ./evm-runners-levels
2. Creating a .env file in ~/.config/evm-runners/`,

	RunE: func(cmd *cobra.Command, args []string) error {
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

		// Ask user if they want to init evm-runners in the current directory
		fmt.Printf("Initializing evm-runners in '%s'\n", subdir)
		fmt.Printf("Continue? (y/n): ")
		var overwrite string
		_, err = fmt.Scanln(&overwrite)
		if err != nil {
			return fmt.Errorf("error reading input: %w", err)
		}

		// print new line
		fmt.Printf("\n")

		if overwrite != "y" && overwrite != "Y" {
			fmt.Printf("Aborting initialization\n")
			return nil
		}

		// Clone ethernautdao/evm-runners-levels.git
		fmt.Printf("Cloning 'github.com/ethernautdao/evm-runners-levels.git' ...\n")

		err = cloneRepository(subdir)
		if err != nil {
			return err
		}

		// Create or update .env file
		envDirPath := filepath.Join(usr.HomeDir, ".evm-runners")
		envFilePath := filepath.Join(envDirPath, ".env")

		fmt.Printf("\nUpdating .env file at '%s' ...\n", envFilePath)

		err = createOrUpdateEnv(subdir, envDirPath, envFilePath)
		if err != nil {
			return err
		}

		fmt.Println("\nevm-runners initialized successfully!\nRun 'evm-runners start' to begin solving a level!")
		return nil
	},
}

func cloneRepository(subdir string) error {
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		execCmd := exec.Command("git", "clone", "https://github.com/ethernautdao/evm-runners-levels.git", subdir)
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		if err := execCmd.Run(); err != nil {
			return fmt.Errorf("error cloning the repository: %v", err)
		}
		fmt.Println("evm-runners-levels cloned successfully")
	} else {
		fmt.Println("evm-runners-levels already exists")
	}

	return nil
}

func createOrUpdateEnv(subdir string, envDirPath string, envFilePath string) error {

	// Set the fields in the config struct
	config := utils.Config{
		EVMR_SERVER:     "https://api.evmr.sh/",
		EVMR_LEVELS_DIR: subdir,
	}

	// Check if the .env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist
		if _, err := os.Stat(envDirPath); os.IsNotExist(err) {
			if err := os.MkdirAll(envDirPath, 0755); err != nil {
				return fmt.Errorf("error creating directory %s: %v", envDirPath, err)
			}
		}

		// create the .env file
		f, err := os.Create(envFilePath)
		if err != nil {
			return fmt.Errorf("error creating .env file: %v", err)
		}
		defer f.Close()

		// write the config to the .env file
		err = utils.WriteConfig(config)
		if err != nil {
			return fmt.Errorf("error writing to .env file: %v", err)
		}

		fmt.Println(".env file created successfully.")
	} else {
		// .env file already exists

		// Load existing config
		existingConfig, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		// set server and levels_dir
		existingConfig.EVMR_SERVER = config.EVMR_SERVER
		existingConfig.EVMR_LEVELS_DIR = config.EVMR_LEVELS_DIR

		// Update config
		err = utils.WriteConfig(existingConfig)
		if err != nil {
			return fmt.Errorf("error writing to .env file: %v", err)
		}

		fmt.Println(".env file updated successfully.")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
