package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display active account and authentication state",
	Run: func(cmd *cobra.Command, args []string) {
		usr, _ := user.Current()
		tokenFile := filepath.Join(usr.HomeDir, ".config", "bb", "token.json")
		data, err := os.ReadFile(tokenFile)
		if err != nil {
			fmt.Println("Not logged in.")
			return
		}

		var tokenData map[string]interface{}
		if err := json.Unmarshal(data, &tokenData); err != nil {
			fmt.Println("Error reading authentication state:", err)
			return
		}

		fmt.Printf("Logged in with token: %s\n", tokenData["access_token"])
	},
}

func init() {
	authCmd.AddCommand(authStatusCmd)
}
