package deployment

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"pushnpray/internal/manifest"
)

func RunDeployment(dep models.Deployment, project models.Project, strategy GitFetchStrategy, manifestFile string) {
	updateStatus := func(status, message, url string) {
		result := database.GetDB().Model(&dep).Updates(models.Deployment{Status: status, Message: message, URL: url})
		if result.Error != nil {
			log.Printf("failed to update deployment %s status: %v", dep.ID, result.Error)
		}
	}

	workspaceDir := filepath.Join("/tmp", "pushnpray-deployments", dep.ID)
	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		updateStatus(StatusError, fmt.Sprintf("%s: %v", msgWorkspaceFailed, err), "")
		return
	}
	defer func() {
		var _ = os.RemoveAll(workspaceDir)
	}()

	fmt.Printf("Fetching repo %s into %s...\n", project.RepositoryUrl, workspaceDir)
	if err := strategy.Fetch(project.RepositoryUrl, workspaceDir); err != nil {
		updateStatus(StatusError, fmt.Sprintf("%s: %v", msgFetchFailed, err), "")
		return
	}

	manifestPath := filepath.Join(workspaceDir, manifestFile)
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		updateStatus(StatusError, fmt.Sprintf(msgManifestMissing, manifestFile), "")
		return
	}

	projectConfig, err := manifest.Unmarshal(manifestPath)
	if err != nil {
		updateStatus(StatusError, fmt.Sprintf("%s: %v", msgManifestInvalid, err), "")
		return
	}

	deployService := NewDeployService()
	if err := deployService.DeployProject(project.Slug, project.ID, projectConfig, workspaceDir); err != nil {
		updateStatus(StatusError, fmt.Sprintf("%s: %v", msgDeployFailed, err), "")
		return
	}

	url := fmt.Sprintf("https://%s-%s.pushnpray.polydo.dev", project.Slug, project.ID)
	updateStatus(StatusSuccess, msgDeploySuccess, url)
}
