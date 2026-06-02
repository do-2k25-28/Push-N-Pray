package routes

import (
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"pushnpray/internal"

	"github.com/gin-gonic/gin"
)

const (
	ErrServiceNotFound     = "Service not found"
	ErrServiceListFailed   = "Failed to list services"
	ErrServiceDeleteFailed = "Failed to delete service"
)

func ListServices(c *gin.Context) {
	project := c.MustGet("project").(models.Project)

	var services []models.ManagedService
	if err := database.GetDB().Where("project_id = ?", project.ID).Order("created_at desc").Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrServiceListFailed})
		return
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}

func DeleteService(c *gin.Context) {
	project := c.MustGet("project").(models.Project)
	serviceID := c.Param("serviceId")

	var service models.ManagedService
	if err := database.GetDB().First(&service, "id = ? AND project_id = ?", serviceID, project.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrServiceNotFound})
		return
	}

	dockerClient, err := internal.NewClient(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrServiceDeleteFailed})
		return
	}

	if err := dockerClient.RemoveManagedServiceResources(c.Request.Context(), service.ContainerName, service.VolumeName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrServiceDeleteFailed})
		return
	}

	service.Status = models.ServiceDeleted
	if err := database.GetDB().Save(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrServiceDeleteFailed})
		return
	}

	c.Status(http.StatusNoContent)
}
