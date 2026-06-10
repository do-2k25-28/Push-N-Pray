package routes

import (
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetProjectEnv(c *gin.Context) {
	project := c.MustGet("project").(models.Project)

	var records []models.EnvVar
	if result := database.GetDB().Where("project = ?", project.ID).Find(&records); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve environment variables"})
		return
	}

	type envVarResponse struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	variables := make([]envVarResponse, 0, len(records))
	for _, r := range records {
		variables = append(variables, envVarResponse{Name: r.Name, Value: r.Value})
	}

	c.JSON(http.StatusOK, gin.H{"variables": variables})
}

func DeleteProjectEnv(c *gin.Context) {
	project := c.MustGet("project").(models.Project)
	name := c.Param("name")

	result := database.GetDB().
		Where("project = ? AND name = ?", project.ID, name).
		Delete(&models.EnvVar{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete environment variable"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "environment variable not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

type setEnvRequest struct {
	Variables []struct {
		Name  string `json:"name"  binding:"required"`
		Value string `json:"value"`
	} `json:"variables" binding:"required"`
}

func SetProjectEnv(c *gin.Context) {
	project := c.MustGet("project").(models.Project)

	var req setEnvRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	records := make([]models.EnvVar, 0, len(req.Variables))
	for _, v := range req.Variables {
		records = append(records, models.EnvVar{
			Project: project.ID,
			Name:    v.Name,
			Value:   v.Value,
		})
	}

	if len(records) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	result := database.GetDB().
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(&records)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save environment variables"})
		return
	}

	c.Status(http.StatusNoContent)
}
