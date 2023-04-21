package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os/exec"
)

type AuthResponse struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Discriminator int    `json:"discriminator"`
	AccessToken   string `json:"access_token"`
}

type Config struct {
	EVMR_SERVER string `mapstructure:"EVMR_SERVER"`
	EVMR_TOKEN  string `mapstructure:"EVMR_TOKEN"`
	EVMR_ID     string `mapstructure:"EVMR_ID"`
	EVMR_NAME   string `mapstructure:"EVMR_NAME"`
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

			tokenUrl := fmt.Sprintf("https://evm-runners.fly.dev/users/info/%s", pin)
			resp, err := http.Get(tokenUrl)
			if err != nil {
				return fmt.Errorf("error making GET request: %v", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("error reading response body: %v", err)
			}

			var authResp AuthResponse
			if err := json.Unmarshal(body, &authResp); err != nil {
				return fmt.Errorf("error unmarshalling response body: %v", err)
			}

			if err := saveDataToEnv(authResp); err != nil {
				return fmt.Errorf("failed to save auth data: %v", err)
			}
		} else {
			return fmt.Errorf("Invalid authentication method. Only Discord is available yet.")
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

func saveDataToEnv(authResp AuthResponse) error {
	/*     homeDir, err := os.UserHomeDir()
	    if err != nil {
	        return fmt.Errorf("failed to get home directory: %v", err)
	    }

	    configFile := fmt.Sprintf("%s/.env", homeDir)
		    viper.SetConfigFile(configFile)
	*/

	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	viper.Set("EVMR_TOKEN", authResp.AccessToken)
	viper.Set("EVMR_ID", authResp.ID)
	viper.Set("EVMR_NAME", fmt.Sprintf("%s#%04d", authResp.Name, authResp.Discriminator))

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %v", err)
	}

	fmt.Println("\nUser data saved in your .env file!")

	return nil
}
