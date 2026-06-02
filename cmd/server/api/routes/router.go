package routes

import (
	"net/http"
	"pushnpray/cmd/server/api/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/v1/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Auth
	auth := router.Group("/v1/auth")

	auth.POST("/register", Register)
	auth.POST("/login", Login)
	auth.POST("/token", Token)

	// PAT
	tokens := router.Group("/v1/tokens")
	router.Use(middleware.Auth())

	tokens.GET("", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	tokens.POST("", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	tokens.DELETE("/:tokenId", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	// Projects
	projects := router.Group("/v1/projects")
	projects.Use(middleware.Auth(), middleware.ProjectOwnership())

	projects.POST("/:projectId", CreateProject)
	projects.GET("", ListProjects)
	projects.GET("/:projectId", GetProject)
	projects.DELETE("/:projectId", DeleteProject)

	// Deployments
	deployments := router.Group("/v1/deployments")
	deployments.Use(middleware.Auth(), middleware.ProjectOwnership())

	deployments.POST("", DeployProject)
	deployments.GET("", ListDeployments)
	deployments.GET("/:deploymentId", GetDeployment)

	// Environment variables
	env := router.Group("/v1/projects/:projectId/env")
	env.Use(middleware.Auth(), middleware.ProjectOwnership())

	env.POST("", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	return router
}
