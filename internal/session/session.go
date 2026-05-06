package session

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// Session represents the authentication state persisted for the CLI.
// The struct is serialized as TOML in the user config file.
type Session struct {
	AccessToken string `toml:"access-token"`
}

// configPath resolves the absolute path of the CLI session file.
// It uses the OS-specific user config directory (for example ~/.config on Linux)
// and appends the Push'N'Pray session file name.
//
// Response example: "/home/ltorres/.config".
func configPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve user config dir: %w", err)
	}

	return filepath.Join(configDir, "pushnpray.toml"), nil
}

// LoadSession reads the session file from disk and decodes its TOML content.
// It returns a Session pointer when the file exists and contains valid data.
// It returns an error if the file cannot be read or parsed.
func LoadSession() (*Session, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read session file: %w", err)
	}

	var session Session
	if err := toml.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &session, nil
}

// VerifyAuth validates that the user is considered logged in.
// It checks that a session file exists and that the access token is not empty.
// It returns a user-facing error message when authentication is missing/invalid.
func VerifyAuth() error {
	session, err := LoadSession()
	if err != nil {
		return errors.New("you are not logged in, run `pushnpray login` first")
	}

	if session.AccessToken == "" {
		return errors.New("session is invalid: missing access token")
	}

	return nil
}
