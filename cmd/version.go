package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Version = "v0.4.1"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of bb-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("bb-cli version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
