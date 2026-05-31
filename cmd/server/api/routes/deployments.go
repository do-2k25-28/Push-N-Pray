package routes

import (
	"net/http"
	"path/filepath"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/deployment"
	"pushnpray/cmd/server/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ErrDeploymentNotFound     = "Deployment not found"
	ErrDeploymentCreateFailed = "Failed to create deployment"
	ErrDeployMissingTarget    = "Must provide tag, commit, or branch"
)

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
	} else if req.Manifest == nil || !filepath.IsAbs(*req.Manifest) {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrDeployMissingTarget})
		return
	}

	deploymentId := uuid.New().String()
	dep := models.Deployment{
		ID:        deploymentId,
		ProjectID: projectId,
		Status:    models.InProgress,
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

	c.JSON(http.StatusAccepted, gin.H{"id": deploymentId})
}
