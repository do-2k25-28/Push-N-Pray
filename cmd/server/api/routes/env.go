package routes

import (
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"pushnpray/cmd/server/utils"

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
		Name   string  `json:"name"`
		Value  *string `json:"value,omitempty"`
		Secret bool    `json:"secret,omitempty"`
	}

	variables := make([]envVarResponse, 0, len(records))
	for _, r := range records {
		v := envVarResponse{Name: r.Name, Secret: r.Secret}
		if !r.Secret {
			val := r.Value
			v.Value = &val
		}
		variables = append(variables, v)
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
		Name   string `json:"name"  binding:"required"`
		Value  string `json:"value"`
		Secret bool   `json:"secret"`
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
	indexByName := make(map[string]int, len(req.Variables))
	for _, v := range req.Variables {
		value := v.Value
		if v.Secret {
			encrypted, err := utils.EncryptSecret(v.Value)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt secret value"})
				return
			}
			value = encrypted
		}
		rec := models.EnvVar{
			Project: project.ID,
			Name:    v.Name,
			Value:   value,
			Secret:  v.Secret,
		}
		if idx, ok := indexByName[v.Name]; ok {
			records[idx] = rec
			continue
		}
		indexByName[v.Name] = len(records)
		records = append(records, rec)
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
