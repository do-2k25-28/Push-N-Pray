package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RouteHealthV1            = "/v1/health"
	RouteAuthRegisterV1      = "/v1/auth/register"
	RouteAuthLoginV1         = "/v1/auth/login"
	RouteAuthTokenV1         = "/v1/auth/token"
	RouteTokensV1            = "/v1/tokens"
	RouteTokenDeleteV1       = "/v1/tokens/:tokenId"
	RouteProjectsV1          = "/v1/projects"
	RouteProjectDeleteV1     = "/v1/projects/:projectId"
	RouteProjectDeployV1     = "/v1/projects/:projectId/deploy"
	RouteProjectDeploymentV1 = "/v1/projects/:projectId/deployments/:deploymentId"
	RouteProjectEnvV1        = "/v1/projects/:projectId/env"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET(RouteHealthV1, func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST(RouteAuthRegisterV1, func(c *gin.Context) {
		c.String(http.StatusOK, RouteAuthRegisterV1)
	})

	r.POST(RouteAuthLoginV1, func(c *gin.Context) {
		c.String(http.StatusOK, RouteAuthLoginV1)
	})

	r.POST(RouteAuthTokenV1, func(c *gin.Context) {
		c.String(http.StatusOK, RouteAuthTokenV1)
	})

	r.GET(RouteTokensV1, func(c *gin.Context) {
		c.String(http.StatusOK, "GET "+RouteTokensV1)
	})

	r.POST(RouteTokensV1, func(c *gin.Context) {
		c.String(http.StatusOK, "POST "+RouteTokensV1)
	})

	r.DELETE(RouteTokenDeleteV1, func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	r.POST(RouteProjectsV1, func(c *gin.Context) {
		c.String(http.StatusOK, RouteProjectsV1)
	})

	r.DELETE(RouteProjectDeleteV1, func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	r.POST(RouteProjectDeployV1, func(c *gin.Context) {
		c.String(http.StatusOK, RouteProjectDeployV1)
	})

	r.GET(RouteProjectDeploymentV1, func(c *gin.Context) {
		c.String(http.StatusOK, RouteProjectDeploymentV1)
	})

	r.POST(RouteProjectEnvV1, func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	return r
}
