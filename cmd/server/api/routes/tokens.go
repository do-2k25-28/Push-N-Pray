package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"pushnpray/cmd/server/database"
	"pushnpray/cmd/server/models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListTokens(c *gin.Context) {
	userID := c.GetString("userID")
	var tokens []models.PersonalAccessToken

	if err := database.GetDB().Where("owner = ?", userID).Find(&tokens).Error; err != nil {
		log.Printf("failed to fetch tokens: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch tokens: %v", err)})
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

	tokenString := strings.ReplaceAll(uuid.New().String(), "-", "")
	hash := sha256.Sum256([]byte(tokenString))
	hashString := hex.EncodeToString(hash[:])

	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		t := time.Unix(*req.ExpiresAt, 0)
		expiresAt = &t
	}

	token := models.PersonalAccessToken{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Hash:      hashString,
		Owner:     userID,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := database.GetDB().Create(&token).Error; err != nil {
		log.Printf("failed to create token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create token: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    token.ID,
		"token": tokenString,
	})
}

func DeleteToken(c *gin.Context) {
	userID := c.GetString("userID")
	tokenID := c.Param("tokenId")

	if err := database.GetDB().Where("id = ? AND owner = ?", tokenID, userID).Delete(&models.PersonalAccessToken{}).Error; err != nil {
		log.Printf("failed to delete token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete token: %v", err)})
		return
	}

	c.Status(http.StatusNoContent)
}
