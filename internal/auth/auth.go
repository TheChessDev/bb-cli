package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
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

// RetrieveCredentials retrieves an API token or OAuth credentials.
func RetrieveCredentials(serverURL string) (string, error) {
	// Step 1: Check for the API token in the token file
	apiToken, err := loadToken()
	if err == nil && apiToken != "" {
		return apiToken, nil
	}

	// Step 2: Fall back to the credentials file
	creds, err := loadCredentials()
	if err != nil {
		return "", errors.New("no valid credentials found, please log in using `bb auth login`")
	}

	if creds.APIToken != "" {
		return creds.APIToken, nil
	}

	return "", errors.New("no valid credentials found")
}

// loadToken reads the API token from the token file
func loadToken() (string, error) {
	tokenPath := getTokenPath()
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New("API token not found")
		}
		return "", fmt.Errorf("failed to read token file: %w", err)
	}
	token := strings.TrimSpace(string(data))
	return token, nil
}

// loadCredentials reads OAuth credentials from the credentials file
func loadCredentials() (*Credentials, error) {
	credsPath := getConfigPath()
	data, err := os.ReadFile(credsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("credentials not found")
		}
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return &creds, nil
}

// getTokenPath returns the path to the API token file
func getTokenPath() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".config", "bb", "token")
}

// getConfigPath returns the path to the credentials file
func getConfigPath() string {
	usr, _ := user.Current()
	return filepath.Join(usr.HomeDir, ".config", "bb", "credentials.json")
}
