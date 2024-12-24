package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/TheChessDev/bb-cli/internal/auth"
	"github.com/spf13/cobra"
)

type PullRequest struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author struct {
		DisplayName string `json:"display_name"`
	} `json:"author"`
}

type PullRequestResponse struct {
	Values []PullRequest `json:"values"`
}

var prListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pull requests in the current repository",
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Use the validated repository context
		if repositoryContext == "" {
			fmt.Println("Error: repository context not set. Ensure you are in a valid Bitbucket repository.")
			os.Exit(1)
		}

		parts := strings.Split(repositoryContext, "/")
		if len(parts) != 2 {
			fmt.Println("Error: invalid repository context")
			os.Exit(1)
		}
		workspace, repoSlug := parts[0], parts[1]

		// Step 2: Retrieve credentials
		serverURL := "https://api.bitbucket.org"
		creds, err := auth.RetrieveCredentials(serverURL)
		if err != nil {
			fmt.Printf("Error retrieving credentials: %v\n", err)
			os.Exit(1)
		}

		// Step 3: Build the API request
		apiURL := fmt.Sprintf("%s/2.0/repositories/%s/%s/pullrequests", serverURL, workspace, repoSlug)
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}

		// Step 4: Set the Authorization header
		if creds.APIToken != "" {
			req.Header.Set("Authorization", "Bearer "+creds.APIToken)
		} else {
			fmt.Println("No valid credentials found.")
			os.Exit(1)
		}

		// Step 5: Make the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		// Step 6: Handle the response
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Error: received status code %d\n", resp.StatusCode)
			os.Exit(1)
		}

		// Step 7: Decode and display the response
		var prResponse PullRequestResponse
		if err := json.NewDecoder(resp.Body).Decode(&prResponse); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Pull Requests:")
		for _, pr := range prResponse.Values {
			fmt.Printf("- ID: %d | Title: %s | Author: %s\n", pr.ID, pr.Title, pr.Author.DisplayName)
		}
	},
}

func init() {
	prCmd.AddCommand(prListCmd)
}
