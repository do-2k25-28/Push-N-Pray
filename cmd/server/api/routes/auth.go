package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

	if req.RegisterToken != os.Getenv("REGISTER_TOKEN") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid registration token"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	user := models.User{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := database.GetDB().Create(user).Error; err != nil {
		log.Printf("failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	accessToken, err := utils.GetJWTHelper().GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		log.Printf("failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	refreshToken := models.RefreshToken{
		Owner: user.ID,
		Token: uuid.NewString(),
	}

	if err := database.GetDB().Create(refreshToken).Error; err != nil {
		log.Printf("failed to create refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, pkgapi.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})
}
func Login(c *gin.Context) {
	var req pkgapi.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBody})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBody})
		return
	}

	var user models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	match, err := utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil || !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := utils.GetJWTHelper().GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		log.Printf("failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	refreshToken := models.RefreshToken{
		Owner: user.ID,
		Token: uuid.NewString(),
	}

	if err := database.GetDB().Create(refreshToken).Error; err != nil {
		log.Printf("failed to create refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, pkgapi.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
	})
}

func Token(c *gin.Context) {
	var req pkgapi.TokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBody})
		return
	}
	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidBody})
		return
	}

	var rt models.RefreshToken
	if err := database.GetDB().Where("token = ?", req.RefreshToken).First(&rt).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	if err := database.GetDB().Delete(&rt).Error; err != nil {
		log.Printf("failed to delete refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	accessToken, err := utils.GetJWTHelper().GenerateToken(rt.Owner, 24*time.Hour)
	if err != nil {
		log.Printf("failed to generate access token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	newRefreshToken := models.RefreshToken{
		Owner: rt.Owner,
		Token: uuid.NewString(),
	}

	if err := database.GetDB().Create(newRefreshToken).Error; err != nil {
		log.Printf("failed to create refresh token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, pkgapi.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken.Token,
	})
}
