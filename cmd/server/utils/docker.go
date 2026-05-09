package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func BuildDockerImage(imageName, dockerfilePath, context string) error {
	fmt.Printf("Building docker image %s from %s...\n", imageName, dockerfilePath)
	if context == "" {
		context = "."
	}
	cmd := exec.Command("docker", "build", "-t", imageName, "-f", dockerfilePath, context)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", ErrDockerBuildFailed, err)
	}
	return nil
}

func RunDockerContainer(containerName, imageName string) error {
	fmt.Printf("Running docker container %s from image %s...\n", containerName, imageName)
	cmd := exec.Command("docker", "run", "-d", "--name", containerName, imageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %w", ErrDockerRunFailed, err)
	}
	return nil
}

func StopAndRemoveContainersByPattern(pattern string) error {
	out, err := exec.Command("docker", "ps", "-a", "--filter", "name="+pattern, "--format", "{{.Names}}").Output()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDockerListFailed, err)
	}
	for _, name := range strings.Fields(string(out)) {
		var errDockerStop = exec.Command("docker", "stop", name).Run()
		if errDockerStop != nil {
			return errDockerStop
		}

		var errDockerRm = exec.Command("docker", "rm", name).Run()
		if errDockerRm != nil {
			return errDockerRm
		}
	}
	return nil
}

func CheckIfDockerInstalled() bool {
	_, err := exec.LookPath("docker")
	return err == nil
}
