package session

// Config is the root TOML document stored in ~/.config/pushnpray.toml.
// It groups all known sessions by authentication strategy.
type Config struct {
	Sessions Sessions `toml:"sessions"`
}

// Sessions contains the two supported session collections:
// classic for managed Push'N'Pray instances and bearer for self-hosted ones.
type Sessions struct {
	Classic []ClassicSession `toml:"classic"`
	Bearer  []BearerSession  `toml:"bearer"`
}

// ClassicSession stores credentials for the managed instance flow where
// requests are authenticated with Basic auth (email + personal access token).
type ClassicSession struct {
	URL   string `toml:"url"`
	Email string `toml:"email"`
	Token string `toml:"token"`
}

// BearerSession stores credentials for custom/self-hosted instances where
// requests are authenticated with bearer access tokens (+ refresh token).
type BearerSession struct {
	URL          string `toml:"url"`
	AccessToken  string `toml:"access-token"`
	RefreshToken string `toml:"refresh-token"`
}
