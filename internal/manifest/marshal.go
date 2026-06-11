package manifest

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

const DefaultManifestName = "pushnpray.toml"
const defaultServer = "https://pushnpray.hagridshut.net"

func Marshal(manifest Manifest) ([]byte, error) {
	if manifest.Server == defaultServer {
		manifest.Server = ""
	}
	if manifest.UpdateStrategy == UpdateStrategyRecreate {
		manifest.UpdateStrategy = ""
	}

	return toml.Marshal(manifest)
}

func Unmarshal(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Manifest{}, err
	}

	var manifest Manifest
	err = toml.Unmarshal(data, &manifest)
	if err != nil {
		return Manifest{}, err
	}

	if manifest.Server == "" {
		manifest.Server = defaultServer
	}
	if manifest.UpdateStrategy == "" {
		manifest.UpdateStrategy = UpdateStrategyRecreate
	}

	return manifest, nil
}
