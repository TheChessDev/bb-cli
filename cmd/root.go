/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var repositoryContext string

func validateRepositoryContext() error {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	output, err := cmd.Output()
	if err != nil {
		return errors.New("not a Git repository or no remote.origin.url configured")
	}

	remoteURL := strings.TrimSpace(string(output))

	re := regexp.MustCompile(`bitbucket\.org[:/]([\w-]+/[\w-]+)`)
	matches := re.FindStringSubmatch(remoteURL)
	if len(matches) < 2 {
		return errors.New("none of the Git remotes point to a known Bitbucket host")
	}

	repositoryContext = matches[1]

	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bb",
	Short: "Bitbucket CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Parent() != nil && cmd.Parent().Name() == "auth" {
			return nil
		}

		if err := validateRepositoryContext(); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bb-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
