package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"io"
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
	Use:   "auth <discord | wallet>",
	Short: "Authenticate your account or link your wallet address",
	Long: `Authenticate your account or link your wallet address.
Currently, Discord is the only available platform for authentication.

Linking your wallet address enables submissions from the website. To link it, run 
'evmr auth wallet' and enter your wallet address.

If your Discord username has changed and you want to update it, run 
'evmr auth discord' again.`,

	RunE: func(cmd *cobra.Command, args []string) error {

		// commenting this out until more platforms to auth are available
		//if len(args) == 0 {
		//	return fmt.Errorf("Please provide a platform, e.g. Discord\n")
		//}

		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		// discord authentication
		if args[0] == ("discord") || args[0] == ("d") || args[0] == ("Discord") || args[0] == ("") {
			// Check if user authenticated before
			if config.EVMR_TOKEN != "" || config.EVMR_ID != "" || config.EVMR_NAME != "" {
				var overwrite string
				fmt.Printf("It seems like you authenticated before as '%s'\n\nDo you want to update your info? (y/n): ", config.EVMR_NAME)
				fmt.Scanln(&overwrite)
				if overwrite != "y" && overwrite != "Y" {
					fmt.Println("\nAborting authentication")
					return nil
				}
			}

			err := authDiscord(config)
			if err != nil {
				return fmt.Errorf("failed to authenticate with Discord: %v", err)
			}

			fmt.Println("\nSuccessfully authenticated with Discord!")
		} else if args[0] == ("wallet") || args[0] == ("address") {
			fmt.Println("Please enter your wallet address: ")
			var address string
			fmt.Scanln(&address)

			// check if valid ethereum address
			if !utils.IsValidEthereumAddress(address) {
				return fmt.Errorf("Invalid Ethereum address!\n")
			}

			fmt.Printf("\nLinking wallet address %s to your account...\n", address)

			err := linkWallet(config, address)
			if err != nil {
				return fmt.Errorf("failed to link wallet address: %v", err)
			}
		} else {
			return fmt.Errorf("Invalid authentication method!\n")
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

	fmt.Println("When you're done authenticating, enter the provided PIN code")

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

	body, err := io.ReadAll(resp.Body)
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

func linkWallet(config utils.Config, address string) error {
	// Define the JSON payload
	jsonPayload := []byte(fmt.Sprintf(`{"address":"%s"}`, address))

	// Make the HTTP request
	url := config.EVMR_SERVER + "users/wallet"
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.EVMR_TOKEN)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending the request: %v", err)
	}
	defer resp.Body.Close()

	// Check for errors in the response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http request failed with status: %s", resp.Status)
	}

	return nil
}
