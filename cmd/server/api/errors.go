package api

const (
	ErrProjectNotFound        = "Project not found"
	ErrProjectListFailed      = "Failed to list projects"
	ErrProjectCreateFailed    = "Failed to create project"
	ErrProjectDeleteFailed    = "Failed to delete project"
	ErrDeploymentNotFound     = "Deployment not found"
	ErrDeploymentListFailed   = "Failed to list deployments"
	ErrDeploymentCreateFailed = "Failed to create deployment"
	ErrDeployMissingTarget    = "Must provide tag, commit, or branch"
)
