package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type Credentials struct {
	ServerURL    string `json:"server_url"`
	APIToken     string `json:"api_token,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

// SaveAPIToken saves the API token for a specific server URL
func SaveAPIToken(serverURL, apiToken string) error {
	creds := Credentials{
		ServerURL: serverURL,
		APIToken:  apiToken,
	}
	return saveCredentials(creds)
}

// SaveOAuthCredentials saves the OAuth credentials for a specific server URL
func SaveOAuthCredentials(serverURL, clientID, clientSecret string) error {
	creds := Credentials{
		ServerURL:    serverURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
	return saveCredentials(creds)
}

// saveCredentials saves the credentials to a JSON file
func saveCredentials(creds Credentials) error {
	configPath := getConfigPath()
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create credentials file: %w", err)
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(creds)
}

// getConfigPath returns the path to the credentials file
func getConfigPath() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".config", "bb", "credentials.json")
}

// RetrieveCredentials retrieves the credentials for a specific server URL
func RetrieveCredentials(serverURL string) (*Credentials, error) {
	configPath := getConfigPath()

	// Check if the credentials file exists
	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("no credentials found, please log in using `bb auth login`")
		}
		return nil, fmt.Errorf("error opening credentials file: %w", err)
	}
	defer file.Close()

	// Decode the credentials file
	var creds Credentials
	if err := json.NewDecoder(file).Decode(&creds); err != nil {
		return nil, fmt.Errorf("error decoding credentials file: %w", err)
	}

	// Verify that the requested server URL matches
	if creds.ServerURL != serverURL {
		return nil, fmt.Errorf("no credentials found for server: %s", serverURL)
	}

	return &creds, nil
}
