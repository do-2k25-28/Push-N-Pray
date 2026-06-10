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
	tokens.Use(middleware.Auth())

	tokens.GET("", ListTokens)
	tokens.POST("", CreateToken)
	tokens.DELETE("/:tokenId", DeleteToken)

	// Projects
	projects := router.Group("/v1/projects")
	projects.Use(middleware.Auth(), middleware.ProjectOwnership())

	projects.POST("", CreateProject)
	projects.GET("", ListProjects)
	projects.GET("/:projectId", GetProject)
	projects.DELETE("/:projectId", DeleteProject)

	// Deployments
	deployments := router.Group("/v1/projects/:projectId/deployments")
	deployments.Use(middleware.Auth(), middleware.ProjectOwnership())

	deployments.POST("", DeployProject)
	deployments.GET("", ListDeployments)
	deployments.GET("/:deploymentId", GetDeployment)

	// Apps
	apps := router.Group("/v1/projects/:projectId/apps")
	apps.Use(middleware.Auth(), middleware.ProjectOwnership())

	apps.GET("/:appName/logs", GetAppLogs)

	// Environment variables
	env := router.Group("/v1/projects/:projectId/env")
	env.Use(middleware.Auth(), middleware.ProjectOwnership())

	env.POST("", SetProjectEnv)

	return router
}
