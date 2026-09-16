package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a password with its hash, supporting bcrypt and legacy SHA256.
func CheckPasswordHash(password, hash string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err == nil {
		return true
	}

	// Legacy SHA256 fallback
	h := sha256.New()
	h.Write([]byte(password))
	shaHash := hex.EncodeToString(h.Sum(nil))
	return subtle.ConstantTimeCompare([]byte(hash), []byte(shaHash)) == 1
}

// EncriptKey maintains backward compatibility for existing SHA-256 callers.
func EncriptKey(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	return hex.EncodeToString(hash.Sum(nil))
}
