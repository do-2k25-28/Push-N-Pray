package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token} or Basic {token}"})
			c.Abort()
			return
		}

		switch parts[0] {
		case "Bearer":
			claims, err := utils.GetJWTHelper().ValidateToken(parts[1])
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
				c.Abort()
				return
			}
			c.Set("userID", claims.UserID)
			c.Next()
			return
		case "Basic":
			// Handle Basic auth for PAT
			email, password, ok := c.Request.BasicAuth()
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid basic auth format"})
				c.Abort()
				return
			}

			// Find user by email
			var user struct {
				ID string
			}
			if err := database.GetDB().Table("users").Select("id").Where("email = ?", email).Scan(&user).Error; err != nil || user.ID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				c.Abort()
				return
			}

			hash := sha256.Sum256([]byte(password))
			hashString := hex.EncodeToString(hash[:])

			// Check PAT
			var token struct {
				ID        string
				ExpiresAt *time.Time
			}
			if err := database.GetDB().Table("personal_access_tokens").Select("id", "expires_at").Where("owner = ? AND hash = ?", user.ID, hashString).Scan(&token).Error; err != nil || token.ID == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid personal access token"})
				c.Abort()
				return
			}

			if token.ExpiresAt != nil && token.ExpiresAt.Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Personal access token expired"})
				c.Abort()
				return
			}

			c.Set("userID", user.ID)
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token} or Basic {token}"})
		c.Abort()
	}
}
