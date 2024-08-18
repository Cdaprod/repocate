package git

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "repocate/internal/utils"
)

// CloneRepository clones a Git repository.
func CloneRepository(workspaceDir, repoURL string) error {
    repoName, err := utils.ExtractRepoName(repoURL)
    if err != nil {
        return fmt.Errorf("failed to extract repo name: %w", err)
    }

    repoPath := utils.GetRepoPath(workspaceDir, repoName)

    cmd := exec.Command("git", "clone", repoURL, repoPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
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
        return err
    }

    cmd = exec.Command("git", "-C", repoPath, "commit", "-m", message)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}