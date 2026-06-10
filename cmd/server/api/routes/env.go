package routes

import (
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

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
