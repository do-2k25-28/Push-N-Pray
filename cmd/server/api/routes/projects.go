package routes

import (
	"fmt"
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"pushnpray/internal"

	pkgapi "pushnpray/pkg/api"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ErrMissingRequiredFields = "slug and repositoryUrl are required"
	ErrProjectNotFound       = "Project not found"
	ErrProjectListFailed     = "Failed to list projects"
	ErrProjectDeleteFailed   = "Failed to delete project"
	ErrDeploymentListFailed  = "Failed to list deployments"
	ErrProjectCreateFailed   = "Failed to create project"
)

func CreateProject(c *gin.Context) {
	userId := c.GetString("userID")

	var req pkgapi.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Slug == "" || req.RepositoryURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrMissingRequiredFields})
		return
	}

	id := uuid.New().String()
	project := models.Project{
		ID:            id,
		Slug:          req.Slug,
		RepositoryUrl: req.RepositoryURL,
		Owner:         userId,
	}

	if err := database.GetDB().Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %v", ErrProjectCreateFailed, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func ListProjects(c *gin.Context) {
	userId := c.GetString("userID")

	var projects []models.Project
	if err := database.GetDB().Find(&projects, "owner = ?", userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrProjectListFailed})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func GetProject(c *gin.Context) {
	project := c.MustGet("project").(models.Project)
	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	project := c.MustGet("project").(models.Project)

	pattern := fmt.Sprintf("%s-%s", project.Slug, project.ID)
	excludedContainers := map[string]struct{}{}
	var services []models.ManagedService
	if err := database.GetDB().Where("project_id = ? AND status <> ?", project.ID, models.ServiceDeleted).Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrServiceListFailed})
		return
	}
	for _, service := range services {
		excludedContainers[service.ContainerName] = struct{}{}
	}

	dockerClient, err := internal.NewClient(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internal.ErrDockerStopRemoveFailed})
		return
	}
	if err := dockerClient.StopContainersByPatternExcept(c.Request.Context(), pattern, excludedContainers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internal.ErrDockerStopRemoveFailed})
		return
	}
	if err := dockerClient.RemoveContainersByPatternExcept(c.Request.Context(), pattern, excludedContainers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internal.ErrDockerStopRemoveFailed})
		return
	}

	if err := database.GetDB().Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrProjectDeleteFailed})
		return
	}
	c.Status(http.StatusNoContent)
}
