package utils

import "errors"

var (
	ErrDockerBuildFailed = errors.New("failed to build docker image")
	ErrDockerRunFailed   = errors.New("failed to run docker container")
	ErrDockerListFailed  = errors.New("failed to list containers")
)
