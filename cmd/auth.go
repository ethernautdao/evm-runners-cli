package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

var authCmd = &cobra.Command{
	Use:   "auth <platform>",
	Short: "Authenticate with Discord",
	Long:  `Authenticate with Discord by visiting a webpage and entering a PIN`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return fmt.Errorf("Please provide a platform, e.g. Discord\n")
		} else if args[0] == ("discord") || args[0] == ("d") || args[0] == ("Discord") {
			url := "https://evm-runners.fly.dev/auth"
			fmt.Printf("Opening %s in your default browser...\n", url)
			if err := openBrowser(url); err != nil {
				return fmt.Errorf("failed to open URL: %v", err)
			}

			fmt.Println("When you're done authenticating, enter the provided PIN code.")

			var pin string
			fmt.Printf("\nPIN: ")
			if _, err := fmt.Scanln(&pin); err != nil {
				return fmt.Errorf("failed to read PIN: %v", err)
			}

			tokenUrl := fmt.Sprintf("https://evm-runners.fly.dev/auth/token/%s", pin)
			fmt.Println("\nRequesting your access token ...")
			resp, err := http.Get(tokenUrl)
			if err != nil {
				return fmt.Errorf("failed to perform GET request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("unexpected status code %d", resp.StatusCode)
			}

			var tokenResp AccessTokenResponse
			if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
				return fmt.Errorf("failed to decode response body: %v", err)
			}

			accessToken := tokenResp.AccessToken
			fmt.Printf("Access token: %s\n", accessToken)

			if err := saveAccessToken(accessToken); err != nil {
				return fmt.Errorf("failed to save access token: %v", err)
			}
		} else {
			fmt.Println("Invalid authentication method. Only Discord is available yet.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func openBrowser(url string) error {
	cmd := exec.Command("open", url) // Use "open" command on macOS
	return cmd.Start()
}

func saveAccessToken(token string) error {
	/*     homeDir, err := os.UserHomeDir()
	    if err != nil {
	        return fmt.Errorf("failed to get home directory: %v", err)
	    }

	    configFile := fmt.Sprintf("%s/.env", homeDir)
		    viper.SetConfigFile(configFile)
	*/

	viper.SetConfigFile(".env")

	viper.Set("EVMR_AUTH", token)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	fmt.Println("\nAccess token saved in your .env file!")

	return nil
}
