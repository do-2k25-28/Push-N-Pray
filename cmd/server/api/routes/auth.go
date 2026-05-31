package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"pushnpray/cmd/server/utils"
	pkgapi "pushnpray/pkg/api"
)

var (
	ErrInvalidBody = "Invalid body"
)

func Register(c *gin.Context) {
	var req pkgapi.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBody})
		return
	}
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBody})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	user := models.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := database.GetDB().Create(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	accessToken, err := utils.GetJWTHelper().GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	refreshToken := models.RefreshToken{
		Owner: user.ID,
		Token: uuid.NewString(),
	}

	if err := database.GetDB().Create(refreshToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, pkgapi.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})
}

func Login(c *gin.Context) {

}

func Token(c *gin.Context) {

}
