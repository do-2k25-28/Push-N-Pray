package api

import (
	"fmt"
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/deployment"
	"pushnpray/cmd/server/models"
	"pushnpray/cmd/server/utils"
	pkgapi "pushnpray/pkg/api"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeployProjectRequest extends the shared type with the server-only Manifest override.
type DeployProjectRequest struct {
	pkgapi.DeployProjectRequest
	Manifest *string `json:"manifest,omitempty"`
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
	var errStopRemoveContainer = utils.StopAndRemoveContainersByPattern(pattern)
	if errStopRemoveContainer != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrDockerStopRemoveFailed})
	}

	if err := database.GetDB().Delete(&project).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrProjectDeleteFailed})
		return
	}
	c.Status(http.StatusNoContent)
}

func ListDeployments(c *gin.Context) {
	projectId := c.Param("projectId")
	if err := database.GetDB().First(&models.Project{}, "id = ?", projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrProjectNotFound})
		return
	}
	var deployments []models.Deployment
	if err := database.GetDB().Where("project_id = ?", projectId).Order("created_at desc").Find(&deployments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrDeploymentListFailed})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deployments": deployments})
}

func GetDeployment(c *gin.Context) {
	projectId := c.Param("projectId")
	deploymentId := c.Param("deploymentId")
	var dep models.Deployment
	if err := database.GetDB().First(&dep, "id = ? AND project_id = ?", deploymentId, projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrDeploymentNotFound})
		return
	}
	c.JSON(http.StatusOK, dep)
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

func DeployProject(c *gin.Context) {
	projectId := c.Param("projectId")

	var project models.Project
	if err := database.GetDB().First(&project, "id = ?", projectId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrProjectNotFound})
		return
	}

	var req DeployProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var strategy deployment.GitFetchStrategy
	if req.Tag != "" {
		strategy = &deployment.TagStrategy{Tag: req.Tag}
	} else if req.Commit != "" {
		strategy = &deployment.CommitStrategy{Commit: req.Commit}
	} else if req.Branch != "" {
		strategy = &deployment.BranchStrategy{Branch: req.Branch}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrDeployMissingTarget})
		return
	}

	deploymentId := uuid.New().String()
	dep := models.Deployment{
		ID:        deploymentId,
		ProjectID: projectId,
		Status:    "in-progress",
	}

	if err := database.GetDB().Create(&dep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrDeploymentCreateFailed})
		return
	}

	manifestPath := "pushnpray.toml"
	if req.Manifest != nil {
		manifestPath = *req.Manifest
	}

	go deployment.RunDeployment(dep, project, strategy, manifestPath)

	c.JSON(http.StatusOK, gin.H{"id": deploymentId})
}
