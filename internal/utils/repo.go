package utils

import (
    "net/url"
    "path/filepath"
    "strings"
    "os"
)

// ExtractRepoName extracts the repository name from the repository URL
func ExtractRepoName(repoURL string) (string, error) {
    parsedURL, err := url.Parse(repoURL)
    if err != nil {
        return "", err
    }

    repoName := filepath.Base(parsedURL.Path) // Use filepath.Base instead of path.Base
    repoName = strings.TrimSuffix(repoName, ".git")

    return repoName, nil
}

// GetRepoPath returns the path of the repository in the workspace
func GetRepoPath(workspaceDir, repoName string) string {
    return filepath.Join(workspaceDir, repoName) // Use filepath.Join instead of path.Join
}

// IsRepoCloned checks if the repository has already been cloned in the workspace
func IsRepoCloned(workspaceDir, repoName string) bool {
    repoPath := GetRepoPath(workspaceDir, repoName)
    if _, err := os.Stat(filepath.Join(repoPath, ".git")); !os.IsNotExist(err) { // Use filepath.Join
        return true
    }
    return false
}