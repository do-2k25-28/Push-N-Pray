package deployment

import (
	"fmt"
	"os"
	"os/exec"
)

type fetchProjectStrategyType int

const (
	fetchStrategyGit fetchProjectStrategyType = iota
	fetchStrategyDockerfile
	fetchStrategyDockerImage
)

var stateName = map[fetchProjectStrategyType]string{
	fetchStrategyGit:         "git strategy",
	fetchStrategyDockerfile:  "dockerfile strategy",
	fetchStrategyDockerImage: "docker image strategy",
}

type fetchProjectStrategyService interface {
	fetch(objectReference string) error
}

type fetchProjectStrategyGitService struct{}

func (s *fetchProjectStrategyGitService) fetch(objectReference string) error {
	cmd := exec.Command("git", "clone", objectReference)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type fetchProjectStrategyDockerImageService struct{}

func (s *fetchProjectStrategyDockerImageService) fetch(objectReference string) error {
	cmd := exec.Command("docker", "pull", objectReference)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type fetchProjectStrategyDockerfileService struct{}

func (s *fetchProjectStrategyDockerfileService) fetch(objectReference string) error {
	_, err := os.Stat(objectReference)
	if os.IsNotExist(err) {
		return fmt.Errorf(errFmtDockerfileNotFound, objectReference)
	}
	return err
}

type fetchProjectService struct {
	strategy fetchProjectStrategyType
	service  fetchProjectStrategyService
}

func NewFetchProjectService(strategy fetchProjectStrategyType) *fetchProjectService {
	var service fetchProjectStrategyService

	switch strategy {
	case fetchStrategyGit:
		service = &fetchProjectStrategyGitService{}
	case fetchStrategyDockerImage:
		service = &fetchProjectStrategyDockerImageService{}
	case fetchStrategyDockerfile:
		service = &fetchProjectStrategyDockerfileService{}
	}

	return &fetchProjectService{
		strategy: strategy,
		service:  service,
	}
}

func (fps *fetchProjectService) getStrategy() fetchProjectStrategyType {
	return fps.strategy
}

func (fps *fetchProjectService) fetch(objectReference string) error {
	if fps.service == nil {
		return ErrNoStrategyConfigured
	}
	return fps.service.fetch(objectReference)
}
