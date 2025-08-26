package persistence

import (
	"os"
	"path/filepath"
)

func ResolveHomeDirPath(path string) string {
	if len(path) > 0 && path[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path // Fallback to original path if home directory cannot be resolved
		}
		return filepath.Join(homeDir, path[1:])
	}
	return path
}

func FileAtHomeDir(parts ...string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(parts...) // Fallback to original path if home directory cannot be resolved
	}
	allParts := append([]string{homeDir}, parts...)
	return filepath.Join(allParts...)
}
