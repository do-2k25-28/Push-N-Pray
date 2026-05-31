package routes

import (
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/deployment"
	"pushnpray/cmd/server/models"
	pkgapi "pushnpray/pkg/api"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ErrDeploymentNotFound     = "Deployment not found"
	ErrDeploymentCreateFailed = "Failed to create deployment"
	ErrDeployMissingTarget    = "Must provide tag, commit, or branch"
)

func ListDeployments(c *gin.Context) {
	project := c.MustGet("project").(models.Project)

	var deployments []models.Deployment
	if err := database.GetDB().Where("project_id = ?", project.ID).Order("created_at desc").Find(&deployments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrDeploymentListFailed})
		return
	}
	c.JSON(http.StatusOK, gin.H{"deployments": deployments})
}

func GetDeployment(c *gin.Context) {
	project := c.MustGet("project").(models.Project)

	deploymentId := c.Param("deploymentId")
	var dep models.Deployment
	if err := database.GetDB().First(&dep, "id = ? AND project_id = ?", deploymentId, project.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrDeploymentNotFound})
		return
	}
	c.JSON(http.StatusOK, dep)
}

func DeployProject(c *gin.Context) {
	project := c.MustGet("project").(models.Project)

	var req pkgapi.DeployProjectRequest
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown deployment strategy"})
		return
	}

	deploymentId := uuid.New().String()
	dep := models.Deployment{
		ID:        deploymentId,
		ProjectID: project.ID,
		Status:    models.InProgress,
	}

	if err := database.GetDB().Create(&dep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrDeploymentCreateFailed})
		return
	}

	go deployment.RunDeployment(dep, project, strategy)

	c.JSON(http.StatusAccepted, gin.H{"id": deploymentId})
}
