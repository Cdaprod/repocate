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

    repoName := filepath.Base(parsedURL.Path)
    repoName = strings.TrimSuffix(repoName, ".git")

    return repoName, nil
}

// GetRepoPath returns the path of the repository in the workspace
func GetRepoPath(workspaceDir, repoName string) string {
    return filepath.Join(workspaceDir, repoName)
}

// IsRepoCloned checks if the repository has already been cloned in the workspace
func IsRepoCloned(workspaceDir, repoName string) bool {
    repoPath := GetRepoPath(workspaceDir, repoName)
    if _, err := os.Stat(filepath.Join(repoPath, ".git")); !os.IsNotExist(err) {
        return true
    }
    return false
}

// CloneRepository clones a Git repository.
func CloneRepository(workspaceDir, repoURL string) error {
    repoName, err := ExtractRepoName(repoURL)
    if err != nil {
        log.Error(fmt.Sprintf("Failed to extract repo name: %s", err))
        return fmt.Errorf("failed to extract repo name: %w", err)
    }

    repoPath := GetRepoPath(workspaceDir, repoName)

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
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "checkout", "-b", branchName)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to create branch %s: %s", branchName, err))
        return err
    }

    log.Info(fmt.Sprintf("Branch %s created successfully.", branchName))
    return nil
}

// CommitChanges commits changes to the repository.
func CommitChanges(workspaceDir, repoName, message string) error {
    repoPath := GetRepoPath(workspaceDir, repoName)

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

// CheckoutBranch checks out an existing branch in the repository.
func CheckoutBranch(workspaceDir, repoName, branchName string) error {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "checkout", branchName)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to checkout branch %s: %s", branchName, err))
        return err
    }

    log.Info(fmt.Sprintf("Checked out branch %s successfully.", branchName))
    return nil
}

// PullChanges pulls the latest changes from the remote repository.
func PullChanges(workspaceDir, repoName string) error {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "pull")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to pull changes: %s", err))
        return err
    }

    log.Info("Changes pulled successfully.")
    return nil
}

// PushChanges pushes local commits to the remote repository.
func PushChanges(workspaceDir, repoName string) error {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "push")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to push changes: %s", err))
        return err
    }

    log.Info("Changes pushed successfully.")
    return nil
}

// ListBranches lists all branches in the repository.
func ListBranches(workspaceDir, repoName string) ([]string, error) {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "branch", "--format=%(refname:short)")
    output, err := cmd.Output()
    if err != nil {
        log.Error(fmt.Sprintf("Failed to list branches: %s", err))
        return nil, err
    }

    branches := strings.Split(strings.TrimSpace(string(output)), "\n")
    return branches, nil
}

// GetCurrentBranch gets the name of the current branch.
func GetCurrentBranch(workspaceDir, repoName string) (string, error) {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
    output, err := cmd.Output()
    if err != nil {
        log.Error(fmt.Sprintf("Failed to get current branch: %s", err))
        return "", err
    }

    return strings.TrimSpace(string(output)), nil
}

// IsRepoClean checks if the repository has any uncommitted changes.
func IsRepoClean(workspaceDir, repoName string) (bool, error) {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
    output, err := cmd.Output()
    if err != nil {
        log.Error(fmt.Sprintf("Failed to check repo status: %s", err))
        return false, err
    }

    return len(output) == 0, nil
}

// CreateSnapshot creates a Git snapshot (commit) of the current state.
func CreateSnapshot(workspaceDir, repoName, message string) error {
    clean, err := IsRepoClean(workspaceDir, repoName)
    if err != nil {
        return err
    }

    if clean {
        log.Info("Repository is clean. No changes to snapshot.")
        return nil
    }

    return CommitChanges(workspaceDir, repoName, message)
}

// Rollback reverts the repository to the last known good commit.
func Rollback(workspaceDir, repoName string) error {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "reset", "--hard", "HEAD^")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to rollback: %s", err))
        return err
    }

    log.Info("Rolled back to the previous commit successfully.")
    return nil
}

// GetRemoteURL gets the URL of the remote repository.
func GetRemoteURL(workspaceDir, repoName string) (string, error) {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "config", "--get", "remote.origin.url")
    output, err := cmd.Output()
    if err != nil {
        log.Error(fmt.Sprintf("Failed to get remote URL: %s", err))
        return "", err
    }

    return strings.TrimSpace(string(output)), nil
}

// SetRemoteURL sets the URL of the remote repository.
func SetRemoteURL(workspaceDir, repoName, url string) error {
    repoPath := GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "-C", repoPath, "remote", "set-url", "origin", url)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        log.Error(fmt.Sprintf("Failed to set remote URL: %s", err))
        return err
    }

    log.Info(fmt.Sprintf("Remote URL set to %s", url))
    return nil
}