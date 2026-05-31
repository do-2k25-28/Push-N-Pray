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

// DeployProjectRequest extends the shared type with the server-only Manifest override.
type DeployProjectRequest struct {
	pkgapi.DeployProjectRequest
	Manifest *string `json:"manifest,omitempty"`
}

func CreateProject(c *gin.Context) {
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
	}

	if err := database.GetDB().Create(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %v", ErrProjectCreateFailed, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func ListProjects(c *gin.Context) {
	var projects []models.Project
	if err := database.GetDB().Find(&projects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrProjectListFailed})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func GetProject(c *gin.Context) {
	projectId := c.Param("projectId")
	var project models.Project
	if err := database.GetDB().First(&project, "id = ?", projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrProjectNotFound})
		return
	}
	c.JSON(http.StatusOK, project)
}

func DeleteProject(c *gin.Context) {
	projectId := c.Param("projectId")
	var project models.Project
	if err := database.GetDB().First(&project, "id = ?", projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrProjectNotFound})
		return
	}

	pattern := fmt.Sprintf("%s-%s", project.Slug, project.ID)
	dockerClient, err := internal.NewClient(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internal.ErrDockerStopRemoveFailed})
		return
	}
	if err := dockerClient.StopContainersByPattern(c.Request.Context(), pattern); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internal.ErrDockerStopRemoveFailed})
		return
	}
	if err := dockerClient.RemoveContainersByPattern(c.Request.Context(), pattern); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": internal.ErrDockerStopRemoveFailed})
	}

	if err := database.GetDB().Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrProjectDeleteFailed})
		return
	}
	c.Status(http.StatusNoContent)
}
