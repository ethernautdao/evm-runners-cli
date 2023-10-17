package cmd

import (
	"bytes"
	"fmt"
	"github.com/ethernautdao/evm-runners-cli/internal/utils"
	"github.com/spf13/cobra"
	"net/http"
)

// addressCmd represents the address command
var addressCmd = &cobra.Command{
	Use:   "address [address]",
	Short: "Link your Optimism address to your account",
	Long: `Link your Optimism address to your account.
Linking your wallet address enables submissions from the website.

Additionally, you will receive NFTs on Optimism for completing levels,
showcasing your skills.

To update your wallet address, simply run this command again.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// load config
		config, err := utils.LoadConfig()
		if err != nil {
			return err
		}

		fmt.Println("Please enter your Optimism address: ")

		// declare address
		var address string

		if len(args) != 0 {
			address = args[0]
		} else {
			fmt.Scanln(&address)
		}

		// check if valid ethereum address
		if !utils.IsValidEthereumAddress(address) {
			return fmt.Errorf("Invalid Optimism address!\n")
		}

		fmt.Printf("\nLinking Optimism address '%s' to your account...\n\n", address)

		err = linkWallet(config, address)
		if err != nil {
			return fmt.Errorf("failed to link wallet address: %v", err)
		} else {
			fmt.Printf("Success!\n")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addressCmd)
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
