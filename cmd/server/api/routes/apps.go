package routes

import (
	"fmt"
	"net/http"
	"pushnpray/cmd/server/models"
	"pushnpray/internal/dockerw"

	"github.com/gin-gonic/gin"
	"github.com/moby/moby/api/pkg/stdcopy"
)

const (
	ErrAppNotFound   = "App not found"
	ErrDockerConnect = "Failed to connect to Docker"
)

func GetAppLogs(c *gin.Context) {
	project := c.MustGet("project").(models.Project)
	appName := c.Param("appName")
	tail := c.DefaultQuery("tail", "100")
	follow := c.Query("follow") == "1"

	dockerClient, err := dockerw.NewClient(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrDockerConnect})
		return
	}

	containerName := fmt.Sprintf("app-%s-%s", project.ID, appName)
	reader, err := dockerClient.GetContainerLogs(c.Request.Context(), containerName, tail, follow)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": ErrAppNotFound})
		return
	}
	defer reader.Close()

	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.Status(http.StatusOK)

	fw := &flushWriter{w: c.Writer}
	if flusher, ok := c.Writer.(http.Flusher); ok {
		fw.flusher = flusher
	}

	stdcopy.StdCopy(fw, fw, reader)
}

type flushWriter struct {
	w       gin.ResponseWriter
	flusher http.Flusher
}

func (fw *flushWriter) Write(p []byte) (int, error) {
	n, err := fw.w.Write(p)
	if err == nil && fw.flusher != nil {
		fw.flusher.Flush()
	}
	return n, err
}
