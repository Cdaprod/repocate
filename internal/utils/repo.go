package utils

import (
    "net/url"
    "path"
    "strings"
    "os"
)

// ExtractRepoName extracts the repository name from the repository URL
func ExtractRepoName(repoURL string) (string, error) {
    parsedURL, err := url.Parse(repoURL)
    if err != nil {
        return "", err
    }

    repoName := path.Base(parsedURL.Path)
    repoName = strings.TrimSuffix(repoName, ".git")

    return repoName, nil
}

// GetRepoPath returns the path of the repository in the workspace
func GetRepoPath(workspaceDir, repoName string) string {
    return path.Join(workspaceDir, repoName)
}

// IsRepoCloned checks if the repository has already been cloned in the workspace
func IsRepoCloned(workspaceDir, repoName string) bool {
    repoPath := GetRepoPath(workspaceDir, repoName)
    if _, err := os.Stat(path.Join(repoPath, ".git")); !os.IsNotExist(err) {
        return true
    }
    return false
}