package session

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"pushnpray/pkg/api"

	"github.com/pelletier/go-toml/v2"
)

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

// Load reads the session file from disk and decodes its TOML content.
// It returns an empty config when the file does not exist yet.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("read session file: %w", err)
	}

	var sessionConfig Config
	if err := toml.Unmarshal(data, &sessionConfig); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &sessionConfig, nil
}

// Save writes the complete session config to disk.
func Save(sessionConfig *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := toml.Marshal(sessionConfig)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write session file: %w", err)
	}

	return nil
}

func SaveClassicSession(url, email, token string) error {
	sessionConfig, err := Load()
	if err != nil {
		return err
	}

	for i := range sessionConfig.Sessions.Classic {
		if sessionConfig.Sessions.Classic[i].URL == url {
			sessionConfig.Sessions.Classic[i].URL = url
			sessionConfig.Sessions.Classic[i].Email = email
			sessionConfig.Sessions.Classic[i].Token = token
			return Save(sessionConfig)
		}
	}

	sessionConfig.Sessions.Classic = append(sessionConfig.Sessions.Classic, ClassicSession{
		URL:   url,
		Email: email,
		Token: token,
	})

	return Save(sessionConfig)
}

func SaveBearerSession(url, accessToken, refreshToken string) error {
	sessionConfig, err := Load()
	if err != nil {
		return err
	}

	for i := range sessionConfig.Sessions.Bearer {
		if sessionConfig.Sessions.Bearer[i].URL == url {
			sessionConfig.Sessions.Bearer[i].URL = url
			sessionConfig.Sessions.Bearer[i].AccessToken = accessToken
			sessionConfig.Sessions.Bearer[i].RefreshToken = refreshToken
			return Save(sessionConfig)
		}
	}

	sessionConfig.Sessions.Bearer = append(sessionConfig.Sessions.Bearer, BearerSession{
		URL:          url,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	return Save(sessionConfig)
}

// GetAuthClientOption returns an API authentication option for the given server URL.
// It picks credentials from the matching saved session (classic first, then bearer).
func GetAuthClientOption(serverURL string) (api.Option, error) {
	if serverURL == "" {
		return nil, fmt.Errorf("server URL is required")
	}

	sessionConfig, err := Load()
	if err != nil {
		return nil, err
	}

	for _, s := range sessionConfig.Sessions.Classic {
		if s.URL == serverURL && s.Email != "" && s.Token != "" {
			return api.WithBasicAuth(s.Email, s.Token), nil
		}
	}

	for _, s := range sessionConfig.Sessions.Bearer {
		if s.URL == serverURL && s.AccessToken != "" {
			return api.WithBearerToken(s.AccessToken), nil
		}
	}

	return nil, fmt.Errorf("no valid session found for this server, run `pushnpray login --server <url>` first")
}

// VerifyAuth validates that the user is considered logged in.
// It checks that a session file exists and that the access token is not empty.
// It returns a user-facing error message when authentication is missing/invalid.
func VerifyAuth() error {
	sessionConfig, err := Load()
	if err != nil {
		return fmt.Errorf("you are not logged in, run `pushnpray login` first")
	}

	for _, s := range sessionConfig.Sessions.Classic {
		if s.URL != "" && s.Email != "" && s.Token != "" {
			return nil
		}
	}

	for _, s := range sessionConfig.Sessions.Bearer {
		if s.URL != "" && s.AccessToken != "" {
			return nil
		}
	}

	return fmt.Errorf("you are not logged in, run `pushnpray login` first")
}
