package manifest

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

func Marshal(manifest *Manifest) ([]byte, error) {
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

	return &manifest, nil
}
