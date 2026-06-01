package manifest

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

var defaultServer = "https://pushnpray.polydo.dev"

func Marshal(manifest Manifest) ([]byte, error) {
	if manifest.Server == defaultServer {
		manifest.Server = ""
	}

	return toml.Marshal(manifest)
}

func Unmarshal(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var manifest Manifest
	err = toml.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	if manifest.Server == "" {
		manifest.Server = defaultServer
	}

	return &manifest, nil
}
