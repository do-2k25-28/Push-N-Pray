package middleware

import (
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"

	"github.com/gin-gonic/gin"
)

// This middleware checks if the user owns the project.
// It gets the project id from the path parameter "projectId".
// If no path param, check if ignored.
//
// This middleware must be used in combination of the auth middleware
// otherwise it won't be able to get the user id and will send
// internal server error to the client
func ProjectOwnership() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetString("userID")

		if userId == "" {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}

		projectId := c.Param("projectId")

		// projectId is not a param therefore the middleware should skip validation
		if projectId == "" {
			c.Next()
			c.Abort()
			return
		}

		var project models.Project
		if err := database.GetDB().First(&project, "id = ?", projectId).Error; err != nil {
			c.Status(http.StatusNotFound)
			c.Abort()
			return
		}

		if project.Owner != userId {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("project", project)

		c.Next()
	}
}
