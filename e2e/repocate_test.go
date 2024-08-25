package e2e

import (
    "testing"
    "os/exec"
    "path/filepath"
    "io/ioutil"
    "strings"
    "github.com/cdaprod/repocate/internal/config"
)

func TestRepocateE2E(t *testing.T) {
    // Set up test-specific configuration
    testConfig := config.Config{
        WorkspaceDir: "/tmp/repocate_e2e_test",
        LogFile:      "/tmp/repocate_e2e.log",
        LogLevel:     "debug",
    }
    
    // Ensure clean test environment
    cleanupTestEnvironment(t, testConfig.WorkspaceDir)

    // Run Repocate command
    cmd := exec.Command("repocate", "create", "https://github.com/example/repo.git")
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Repocate command failed: %v\nOutput: %s", err, output)
    }

    // Verify expected outcomes
    assertContainerCreated(t, "example-repo")
    assertRepoCloned(t, testConfig.WorkspaceDir, "repo")
    assertLogContains(t, testConfig.LogFile, "Repository cloned successfully")
}

func cleanupTestEnvironment(t *testing.T, workspaceDir string) {
    // Remove existing workspace and containers
    // ...
}

func assertContainerCreated(t *testing.T, containerName string) {
    cmd := exec.Command("docker", "ps", "-a", "--format", "{{.Names}}")
    output, err := cmd.Output()
    if err != nil {
        t.Fatalf("Failed to list Docker containers: %v", err)
    }
    if !strings.Contains(string(output), containerName) {
        t.Errorf("Expected container %s not found", containerName)
    }
}

func assertRepoCloned(t *testing.T, workspaceDir, repoName string) {
    repoPath := filepath.Join(workspaceDir, repoName)
    if _, err := ioutil.ReadDir(repoPath); err != nil {
        t.Errorf("Repository directory not found: %v", err)
    }
}

func assertLogContains(t *testing.T, logFile, expectedMessage string) {
    logContent, err := ioutil.ReadFile(logFile)
    if err != nil {
        t.Fatalf("Failed to read log file: %v", err)
    }
    if !strings.Contains(string(logContent), expectedMessage) {
        t.Errorf("Expected message not found in logs: %s", expectedMessage)
    }
}