package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

func getCacheDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	var gitRoot string
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			gitRoot = dir
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not a git repository")
		}
		dir = parent
	}

	cacheDir := filepath.Join(gitRoot, ".gct", "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}
	return cacheDir, nil
}

func getCacheKey(prompt string) string {
	hash := sha256.Sum256([]byte(prompt))
	return hex.EncodeToString(hash[:])
}

func readFromCache(key string) (string, bool) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return "", false
	}

	filePath := filepath.Join(cacheDir, key)
	if data, err := os.ReadFile(filePath); err == nil {
		return string(data), true
	}
	return "", false
}

func writeToCache(key string, content string) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return
	}

	filePath := filepath.Join(cacheDir, key)
	_ = os.WriteFile(filePath, []byte(content), 0644)
}
