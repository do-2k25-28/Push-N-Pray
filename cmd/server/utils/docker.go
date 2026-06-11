package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func SetupDockerBad() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
		return
	}

	dockerDir := filepath.Join(homeDir, ".docker")
	configPath := filepath.Join(dockerDir, "config.json")

	err = os.MkdirAll(dockerDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory %s: %v\n", dockerDir, err)
		return
	}

	err = os.WriteFile(configPath, []byte("{}"), 0644)
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}
}
