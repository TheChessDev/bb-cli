package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var authTokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Print the authentication token in use",
	Run: func(cmd *cobra.Command, args []string) {
		usr, _ := user.Current()
		tokenFile := filepath.Join(usr.HomeDir, ".config", "bb", "token.json")
		data, err := os.ReadFile(tokenFile)
		if err != nil {
			fmt.Println("No authentication token found.")
			return
		}

		var tokenData map[string]interface{}
		if err := json.Unmarshal(data, &tokenData); err != nil {
			fmt.Println("Error reading token:", err)
			return
		}

		fmt.Printf("Token: %s\n", tokenData["access_token"])
	},
}

func init() {
	authCmd.AddCommand(authTokenCmd)
}
