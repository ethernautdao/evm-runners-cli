package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os/exec"
	"time"
)

type AuthResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}

type Config struct {
	EVMR_SERVER string `mapstructure:"EVMR_SERVER"`
	EVMR_TOKEN  string `mapstructure:"EVMR_TOKEN"`
	EVMR_ID     string `mapstructure:"EVMR_ID"`
	EVMR_NAME   string `mapstructure:"EVMR_NAME"`
}

var authCmd = &cobra.Command{
	Use:   "auth <platform>",
	Short: "Authenticate your account.",
	Long: `Authenticate your account. Currently Discord is the only available platform. 
For updating your username, run 'evm-runners auth <platform>' again.`,

	RunE: func(cmd *cobra.Command, args []string) error {

		// check if argument is empty
		if len(args) == 0 {
			return fmt.Errorf("Please provide a platform, e.g. Discord\n")
		}

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		// Check if the required environment variables are already set
		if config.EVMR_TOKEN != "" || config.EVMR_ID != "" || config.EVMR_NAME != "" {
			var overwrite string
			fmt.Printf("The following environment variables are already set:\n\nEVMR_TOKEN=%s\nEVMR_ID=%s\nEVMR_NAME=%s\n\nDo you want to update them? (y/n): ", config.EVMR_TOKEN, config.EVMR_ID, config.EVMR_NAME)
			fmt.Scanln(&overwrite)
			if overwrite != "y" && overwrite != "Y" {
				fmt.Println("\nAborting authentication")
				return nil
			}
		}

		// discord authentication
		if args[0] == ("discord") || args[0] == ("d") || args[0] == ("Discord") {
			err := authDiscord(config)
			if err != nil {
				return fmt.Errorf("failed to authenticate with Discord: %v", err)
			}

			fmt.Println("\nSuccessfully authenticated with Discord!")
		} else {
			return fmt.Errorf("Invalid authentication method. Only Discord is available yet.\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func openBrowser(url string) error {
	// wait for one second before opening the url
	time.Sleep(1 * time.Second)

	// open url
	cmd := exec.Command("open", url)
	return cmd.Start()
}

func authDiscord(config utils.Config) error {
	// get URL to open in the browser
	url := config.EVMR_SERVER + "auth"
	fmt.Printf("Opening %s in your default browser...\n", url)
	if err := openBrowser(url); err != nil {
		return fmt.Errorf("failed to open URL: %v", err)
	}

	fmt.Println("When you're done authenticating, enter the provided PIN code.")

	// read PIN from stdin
	var pin string
	fmt.Printf("\nPIN: ")
	if _, err := fmt.Scanln(&pin); err != nil {
		return fmt.Errorf("failed to read PIN: %v", err)
	}

	// make GET request to server
	tokenUrl := fmt.Sprintf("%susers/info/%s", config.EVMR_SERVER, pin)
	resp, err := http.Get(tokenUrl)
	if err != nil {
		return fmt.Errorf("error making GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate with server: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	// unmarshal response body
	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return fmt.Errorf("error unmarshalling response body: %v", err)
	}

	// set or overwrite config elements
	config.EVMR_ID = authResp.ID
	config.EVMR_NAME = fmt.Sprintf("%s", authResp.Name)
	config.EVMR_TOKEN = authResp.AccessToken

	// save config
	if err := utils.WriteConfig(config); err != nil {
		return fmt.Errorf("failed to save auth data: %v", err)
	}

	return nil
}
