package routes

import (
	"fmt"
	"log"
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"pushnpray/internal/dockerw"

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
		log.Printf("%s: %v", ErrProjectCreateFailed, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %v", ErrProjectCreateFailed, err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func ListProjects(c *gin.Context) {
	userId := c.GetString("userID")

	var projects []models.Project
	if err := database.GetDB().Find(&projects, "owner = ?", userId).Error; err != nil {
		log.Printf("%s: %v", ErrProjectListFailed, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %v", ErrProjectListFailed, err)})
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

	dockerClient, err := dockerw.NewClient(c.Request.Context())
	if err != nil {
		log.Printf("failed to create docker client: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to connect to Docker: %v", err)})
		return
	}

	if err := dockerClient.StopContainersByPattern(c.Request.Context(), pattern); err != nil {
		log.Printf("failed to stop containers %s: %v", pattern, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to stop containers: %v", err)})
		return
	}
	if err := dockerClient.RemoveContainersByPattern(c.Request.Context(), pattern); err != nil {
		log.Printf("failed to remove containers %s: %v", pattern, err)
	}

	if err := database.GetDB().Delete(&project).Error; err != nil {
		log.Printf("%s: %v", ErrProjectDeleteFailed, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%s: %v", ErrProjectDeleteFailed, err)})
		return
	}
	c.Status(http.StatusNoContent)
}
