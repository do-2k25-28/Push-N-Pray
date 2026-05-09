package deployment

import "errors"

var (
	ErrNoStrategyConfigured = errors.New("no strategy service configured")
	ErrRepoCloneFailed      = errors.New("failed to clone repository")
	ErrCommitCheckoutFailed = errors.New("failed to checkout commit")
)

const (
	StatusError   = "error"
	StatusSuccess = "success"

	msgWorkspaceFailed = "Failed to create workspace"
	msgFetchFailed     = "Failed to fetch repository"
	msgManifestMissing = "Manifest %s not found in repository"
	msgManifestInvalid = "Failed to parse manifest"
	msgDeployFailed    = "Deployment failed"
	msgDeploySuccess   = "Deployment successful"

	errFmtDockerfileNotFound = "dockerfile does not exist at path: %s"
	errFmtAppDeployFailed    = "failed to deploy app %s"
	errFmtAppRunFailed       = "failed to run app %s"
)
