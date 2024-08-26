package container

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
    "net/url"
    "path/filepath"
    "github.com/cdaprod/repocate/internal/utils"
    "github.com/cdaprod/repocate/internal/log"
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

// CloneRepository clones a Git repository.
func CloneRepository(workspaceDir, repoURL string) error {
    repoName, err := utils.ExtractRepoName(repoURL)
    if err != nil {
        log.Error(fmt.Sprintf("Failed to extract repo name: %s", err))
        return fmt.Errorf("failed to extract repo name: %w", err)
    }

    repoPath := utils.GetRepoPath(workspaceDir, repoName)

    log.Info(fmt.Sprintf("Cloning repository %s into %s", repoURL, repoPath))
    cmd := exec.Command("git", "clone", repoURL, repoPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to clone repository: %s", err))
        return err
    }

    log.Info("Repository cloned successfully.")
    return nil
}

// CreateBranch creates a new branch in the repository.
func CreateBranch(workspaceDir, repoName, branchName string) error {
    repoPath := utils.GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "checkout", "-b", branchName)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

// CommitChanges commits changes to the repository.
func CommitChanges(workspaceDir, repoName, message string) error {
    repoPath := utils.GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "add", ".")
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to add changes: %s", err))
        return err
    }

    cmd = exec.Command("git", "-C", repoPath, "commit", "-m", message)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to commit changes: %s", err))
        return err
    }

    log.Info("Changes committed successfully.")
    return nil
}

