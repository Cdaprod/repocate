package git

import (
    "fmt"
    "os"
    "os/exec"
//    "path/filepath"
    "github.com/cdaprod/repocate/internal/utils"
    "github.com/cdaprod/repocate/internal/log"
)

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