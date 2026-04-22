package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

// GenerateJWTSecret creates a random HS256 secret and writes it to a file.
// If the file already exists, it returns the existing secret.
func GenerateJWTSecret(path string) (string, error) {
	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		return string(data), nil
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random secret: %w", err)
	}

	secret := hex.EncodeToString(bytes)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(path, []byte(secret), 0600); err != nil {
		return "", err
	}
	return secret, nil
}
