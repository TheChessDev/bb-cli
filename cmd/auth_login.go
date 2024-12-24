package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/TheChessDev/bb-cli/internal/auth"
	"github.com/spf13/cobra"
)

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a Bitbucket account",
	Long:  "Authenticate with Bitbucket using API Token or OAuth, and save credentials for future use.",
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Prompt for server type
		var serverType string
		serverOptions := []string{"Bitbucket.org", "Other"}
		serverPrompt := &survey.Select{
			Message: "Where do you want to authenticate?",
			Options: serverOptions,
		}
		err := survey.AskOne(serverPrompt, &serverType)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		var serverURL string
		if serverType == "Other" {
			// Prompt for custom server URL
			err := survey.AskOne(&survey.Input{
				Message: "Enter the custom Bitbucket server URL:",
			}, &serverURL)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		} else {
			serverURL = "https://bitbucket.org"
		}

		// Step 2: Prompt for authentication method
		var authMethod string
		authOptions := []string{"API Token", "OAuth"}
		authPrompt := &survey.Select{
			Message: "Select the authentication method:",
			Options: authOptions,
		}
		err = survey.AskOne(authPrompt, &authMethod)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		// Handle authentication based on the method
		switch authMethod {
		case "API Token":
			handleAPITokenAuth(serverURL)
		case "OAuth":
			handleOAuthAuth(serverURL)
		}
	},
}

// handleAPITokenAuth prompts the user for an API token and saves it
func handleAPITokenAuth(serverURL string) {
	var apiToken string
	err := survey.AskOne(&survey.Password{
		Message: "Enter your API Token:",
	}, &apiToken)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// Save the API token
	if err := auth.SaveAPIToken(serverURL, apiToken); err != nil {
		fmt.Printf("Error saving API token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("API Token saved successfully!")
}

// handleOAuthAuth prompts the user for OAuth client ID and secret and saves them
func handleOAuthAuth(serverURL string) {
	var clientID, clientSecret string
	survey.AskOne(&survey.Input{Message: "Enter your OAuth Client ID:"}, &clientID)
	survey.AskOne(&survey.Password{Message: "Enter your OAuth Client Secret:"}, &clientSecret)

	// Save the OAuth credentials
	if err := auth.SaveOAuthCredentials(serverURL, clientID, clientSecret); err != nil {
		fmt.Printf("Error saving OAuth credentials: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("OAuth credentials saved successfully!")
}

func init() {
	authCmd.AddCommand(authLoginCmd)
}
