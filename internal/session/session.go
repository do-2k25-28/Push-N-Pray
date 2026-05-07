package session

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

	var cfg Config
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &cfg, nil
}

// Save writes the complete session config to disk.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	data, err := toml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write session file: %w", err)
	}

	return nil
}

func SaveClassicSession(url, email, token string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	for i := range cfg.Sessions.Classic {
		if cfg.Sessions.Classic[i].URL == url {
			cfg.Sessions.Classic[i].URL = url
			cfg.Sessions.Classic[i].Email = email
			cfg.Sessions.Classic[i].Token = token
			return Save(cfg)
		}
	}

	cfg.Sessions.Classic = append(cfg.Sessions.Classic, ClassicSession{
		URL:   url,
		Email: email,
		Token: token,
	})

	return Save(cfg)
}

func SaveBearerSession(url, accessToken, refreshToken string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	for i := range cfg.Sessions.Bearer {
		if cfg.Sessions.Bearer[i].URL == url {
			cfg.Sessions.Bearer[i].URL = url
			cfg.Sessions.Bearer[i].AccessToken = accessToken
			cfg.Sessions.Bearer[i].RefreshToken = refreshToken
			return Save(cfg)
		}
	}

	cfg.Sessions.Bearer = append(cfg.Sessions.Bearer, BearerSession{
		URL:          url,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

	return Save(cfg)
}

// VerifyAuth validates that the user is considered logged in.
// It checks that a session file exists and that the access token is not empty.
// It returns a user-facing error message when authentication is missing/invalid.
func VerifyAuth() error {
	cfg, err := Load()
	if err != nil {
		return errors.New("you are not logged in, run `pushnpray login` first")
	}

	for _, s := range cfg.Sessions.Classic {
		if s.URL != "" && s.Email != "" && s.Token != "" {
			return nil
		}
	}

	for _, s := range cfg.Sessions.Bearer {
		if s.URL != "" && s.AccessToken != "" {
			return nil
		}
	}

	return errors.New("you are not logged in, run `pushnpray login` first")
}
