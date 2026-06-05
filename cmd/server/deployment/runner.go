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

func RunDeployment(dep models.Deployment, project models.Project, strategy GitFetchStrategy) {
	reportStatus := func(status models.DeploymentStatus, message string) {
		log.Printf("[%s] [%s] %s", project.ID, status, message)

		result := database.GetDB().Model(&dep).Updates(models.Deployment{Status: status, Message: message})

		if result.Error != nil {
			log.Printf("failed to update deployment %s status: %v", dep.ID, result.Error)

			// Hail mary to notify the user
			database.GetDB().Model(&dep).Updates(models.Deployment{Status: models.Error, Message: "Internal server error"})
		}
	}

	workspaceDir := filepath.Join("/tmp", "pushnpray-deployments", dep.ID)
	manifestPath := filepath.Join(workspaceDir, manifest.DefaultManifestName)

	if err := os.MkdirAll(workspaceDir, 0755); err != nil {
		reportStatus(models.InProgress, fmt.Sprintf("%s: %v", msgWorkspaceFailed, err))
		return
	}
	defer func() {
		if err := os.RemoveAll(workspaceDir); err != nil {
			log.Printf("failed to remove workspace %s: %v", workspaceDir, err)
		}
	}()

	reportStatus(models.InProgress, fmt.Sprintf("Fetching repository %s using strategy %s", project.RepositoryUrl, strategy.DisplayName()))
	if err := strategy.Fetch(project.RepositoryUrl, workspaceDir); err != nil {
		reportStatus(models.Error, fmt.Sprintf("%s: %v", msgFetchFailed, err))
		return
	}

	reportStatus(models.InProgress, "Reading manifest file")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		reportStatus(models.Error, fmt.Sprint(msgManifestMissing))
		return
	}

	man, err := manifest.Unmarshal(manifestPath)
	if err != nil {
		reportStatus(models.Error, fmt.Sprintf("%s: %v", msgManifestInvalid, err))
		return
	}

	reportStatus(models.InProgress, fmt.Sprintf("Updating %d service(s), then deploying %d app(s)", man.GetServiceCount(), man.GetApplicationCount()))
	if err := DeployProject(project.Slug, project.ID, man, workspaceDir); err != nil {
		reportStatus(models.Error, fmt.Sprintf("%s: %v", msgDeployFailed, err))
		return
	}

	url := fmt.Sprintf("https://%s-%s.pushnpray.polydo.dev", project.Slug, project.ID)

	result := database.GetDB().Model(&dep).Updates(models.Deployment{Status: models.Deployed, Message: "Project successfuly deployed", URL: url})
	if result.Error != nil {
		log.Printf("failed to update deployment %s status: %v", dep.ID, result.Error)
	}
}
