package routes

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListTokens(c *gin.Context) {
	userID := c.GetString("userID")
	var tokens []models.PersonalAccessToken

	if err := database.GetDB().Where("owner = ?", userID).Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tokens"})
		return
	}

	type TokenResponse struct {
		ID        string  `json:"id"`
		Name      string  `json:"name"`
		ExpiresAt *string `json:"expiresAt"`
	}

	var response []TokenResponse
	for _, t := range tokens {
		var expiresAt *string
		if t.ExpiresAt != nil {
			val := strconv.FormatInt(t.ExpiresAt.Unix(), 10)
			expiresAt = &val
		}
		response = append(response, TokenResponse{
			ID:        t.ID,
			Name:      t.Name,
			ExpiresAt: expiresAt,
		})
	}
	if response == nil {
		response = make([]TokenResponse, 0)
	}

	c.JSON(http.StatusOK, gin.H{"tokens": response})
}

type CreateTokenRequest struct {
	Name      string `json:"name" binding:"required"`
	ExpiresAt *int64 `json:"expiresAt"`
}

func CreateToken(c *gin.Context) {
	userID := c.GetString("userID")
	var req CreateTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenBytes := uuid.New().String()
	hash := md5.Sum([]byte(tokenBytes))
	tokenHash := hex.EncodeToString(hash[:])

	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t := time.Unix(*req.ExpiresAt, 0)
		expiresAt = &t
	}

	token := models.PersonalAccessToken{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Hash:      tokenHash,
		Owner:     userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := database.GetDB().Create(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    token.ID,
		"token": tokenHash,
	})
}

func DeleteToken(c *gin.Context) {
	userID := c.GetString("userID")
	tokenID := c.Param("tokenId")

	if err := database.GetDB().Where("id = ? AND owner = ?", tokenID, userID).Delete(&models.PersonalAccessToken{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete token"})
		return
	}

	c.Status(http.StatusNoContent)
}
