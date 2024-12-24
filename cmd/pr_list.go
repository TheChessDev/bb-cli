package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/TheChessDev/bb-cli/internal/auth"
	"github.com/spf13/cobra"
)

var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests in the current repository",
	Run: func(cmd *cobra.Command, args []string) {
		// Get server URL (default to Bitbucket.org for now)
		serverURL := "https://bitbucket.org" // Replace with dynamic URL if needed

		// Retrieve credentials
		creds, err := auth.RetrieveCredentials(serverURL)
		if err != nil {
			fmt.Printf("Error retrieving credentials: %v\n", err)
			os.Exit(1)
		}

		// Authenticate using the retrieved credentials
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/rest/api/1.0/pull-requests", serverURL), nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}

		// Use API Token or OAuth credentials
		if creds.APIToken != "" {
			req.Header.Set("Authorization", "Bearer "+creds.APIToken)
		} else if creds.ClientID != "" && creds.ClientSecret != "" {
			// For OAuth, you'd typically use the access token here
			// (OAuth token refresh logic would be implemented separately)
			req.Header.Set("Authorization", "Bearer <ACCESS_TOKEN>")
		} else {
			fmt.Println("No valid credentials found.")
			os.Exit(1)
		}

		// Execute the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		// Handle the response
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: received status code %d\n", resp.StatusCode)
			os.Exit(1)
		}

		fmt.Println("Pull requests fetched successfully!")
		// Decode and display the response (omitted for brevity)
	},
}

func init() {
	prCmd.AddCommand(prListCmd)
}
