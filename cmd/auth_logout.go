package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out of a Bitbucket account",
	Run: func(cmd *cobra.Command, args []string) {
		usr, _ := user.Current()
		tokenFile := filepath.Join(usr.HomeDir, ".config", "bb", "token.json")
		if _, err := os.Stat(tokenFile); os.IsNotExist(err) {
			fmt.Println("No authentication credentials to log out.")
			return
		}

		if err := os.Remove(tokenFile); err != nil {
			fmt.Println("Failed to remove authentication credentials:", err)
			return
		}
		fmt.Println("Logged out successfully.")
	},
}

func init() {
	authCmd.AddCommand(authLogoutCmd)
}
