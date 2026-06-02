package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTTokenHelper handles token generation and verification
type JWTTokenHelper struct {
	secretKey []byte
}

// CustomClaims defines the structure of the data you want to embed in the JWT
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTTokenHelper initializes the helper with a secret key and issuer
func NewJWTTokenHelper(secret string) *JWTTokenHelper {
	return &JWTTokenHelper{
		secretKey: []byte(secret),
	}
}

// GenerateToken creates a signed JWT for a specific user and role, valid for a duration
func (h *JWTTokenHelper) GenerateToken(userID string, duration time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token using the high-level HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	signedToken, err := token.SignedString(h.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken parses and validates the token, returning the custom claims if successful
func (h *JWTTokenHelper) ValidateToken(tokenStr string) (*CustomClaims, error) {
	// Parse the token with our custom claims structure
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect (HMAC / HS256)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return h.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract and return claims if the token is valid
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token claims")
}

var jwtHelperInstance *JWTTokenHelper = nil

func GetJWTHelper() *JWTTokenHelper {
	if jwtHelperInstance != nil {
		return jwtHelperInstance
	}

	jwtHelperInstance = NewJWTTokenHelper("super-secret-passphrase")

	return jwtHelperInstance
}
