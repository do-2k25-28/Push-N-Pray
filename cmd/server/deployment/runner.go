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
	reportStatus := func(status models.DeploymentStatus, message string) {
		log.Printf("[%s] [%s] %s", project.ID, status, message)

		result := database.GetDB().Model(&dep).Updates(models.Deployment{Status: status, Message: message})

		if result.Error != nil {
			log.Printf("failed to update deployment %s status: %v", dep.ID, result.Error)

			// Hail mary to notify the user
			database.GetDB().Model(&dep).Updates(models.Deployment{Status: models.Error, Message: "Internal server error"})
		}
	}

	var manifestPath string
	var workspaceDir string

	if filepath.IsAbs(manifestFile) {
		manifestPath = manifestFile
		workspaceDir = filepath.Dir(manifestFile)
	} else {
		workspaceDir = filepath.Join("/tmp", "pushnpray-deployments", dep.ID)
		if err := os.MkdirAll(workspaceDir, 0755); err != nil {
			reportStatus(models.InProgress, fmt.Sprintf("%s: %v", msgWorkspaceFailed, err))
			return
		}
		defer func() {
			var _ = os.RemoveAll(workspaceDir)
		}()

		reportStatus(models.InProgress, fmt.Sprintf("Fetching repository %s using strategy %s", project.RepositoryUrl, strategy.DisplayName()))
		if err := strategy.Fetch(project.RepositoryUrl, workspaceDir); err != nil {
			reportStatus(models.Error, fmt.Sprintf("%s: %v", msgFetchFailed, err))
			return
		}
		manifestPath = filepath.Join(workspaceDir, manifestFile)
	}

	reportStatus(models.InProgress, fmt.Sprintf("Reading manifest file at %s", manifestFile))
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		reportStatus(models.Error, fmt.Sprintf(msgManifestMissing, manifestFile))
		return
	}

	projectConfig, err := manifest.Unmarshal(manifestPath)
	if err != nil {
		reportStatus(models.Error, fmt.Sprintf("%s: %v", msgManifestInvalid, err))
		return
	}

	deployService := NewDeployService()
	if err := deployService.DeployProject(project.Slug, project.ID, projectConfig, workspaceDir); err != nil {
		reportStatus(models.Error, fmt.Sprintf("%s: %v", msgDeployFailed, err))
		return
	}

	url := fmt.Sprintf("https://%s-%s.pushnpray.polydo.dev", project.Slug, project.ID)

	result := database.GetDB().Model(&dep).Updates(models.Deployment{Status: models.Deployed, Message: "Project successfuly deployed", URL: url})
	if result.Error != nil {
		log.Printf("failed to update deployment %s status: %v", dep.ID, result.Error)
	}
}
