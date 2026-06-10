package api

// AuthResponse holds access and refresh tokens.
type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// RegisterRequest is the payload for account registration.
type RegisterRequest struct {
	Email         string `json:"email"`
	Password      string `json:"password"`
	RegisterToken string `json:"registerToken"`
}

// LoginRequest is the payload for login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenRequest is the payload for refreshing tokens.
type TokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// PAT represents a personal access token metadata entry.
type PAT struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	ExpiresAt string `json:"expiresAt"`
}

// ListPATsResponse wraps a list of PATs.
type ListPATsResponse struct {
	Tokens []PAT `json:"tokens"`
}

// CreatePATRequest is the payload for creating a PAT.
type CreatePATRequest struct {
	Name      string `json:"name"`
	ExpiresAt *int64 `json:"expiresAt"`
}

// CreatePATResponse holds the created PAT and token value.
type CreatePATResponse struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

// CreateProjectRequest is the payload for creating a project.
type CreateProjectRequest struct {
	Slug          string `json:"slug"`
	RepositoryURL string `json:"repositoryUrl"`
}

// CreateProjectResponse holds the created project ID.
type CreateProjectResponse struct {
	ID string `json:"id"`
}

// ListProjectsResponse wraps a list of projects.
type ListProjectsResponse struct {
	Projects []Project `json:"projects"`
}

// Project represents a project instance.
type Project struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	RepositoryURL string `json:"repositoryUrl"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
}

// DeployProjectRequest selects a deploy target.
type DeployProjectRequest struct {
	Tag    string `json:"tag,omitempty"`
	Commit string `json:"commit,omitempty"`
	Branch string `json:"branch,omitempty"`
}

// DeployProjectResponse holds a deployment ID.
type DeployProjectResponse struct {
	ID string `json:"id"`
}

// DeploymentInfoResponse provides deployment status.
type DeploymentInfoResponse struct {
	Status  string `json:"status"`
	URL     string `json:"url,omitempty"`
	Message string `json:"message,omitempty"`
}

// Deployment represents a deployment in the system.
type Deployment struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ListDeploymentsResponse provides a list of deployments.
type ListDeploymentsResponse struct {
	Deployments []Deployment `json:"deployments"`
}

// EnvVar represents a project environment variable in a request.
type EnvVar struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Secret bool   `json:"secret,omitempty"`
}

// ListedEnvVar is returned by GET /env.
// Value is nil for secret variables.
type ListedEnvVar struct {
	Name   string  `json:"name"`
	Value  *string `json:"value,omitempty"`
	Secret bool    `json:"secret,omitempty"`
}

// SetProjectEnvRequest is the payload for setting project env vars.
type SetProjectEnvRequest struct {
	Variables []EnvVar `json:"variables"`
}

// GetProjectEnvResponse wraps the list of env vars for a project.
type GetProjectEnvResponse struct {
	Variables []ListedEnvVar `json:"variables"`
}
