package utils

import (
	"github.com/alexedwards/argon2id"
)

// HashPassword takes a plaintext password and returns the encoded Argon2id string.
func HashPassword(password string) (string, error) {
	// Uses secure default parameters automatically
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

// CheckPasswordHash compares a plaintext password against the stored hash.
func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
