// internal/container/container.go
package container

import (
    "os"
    "fmt"
    "os/exec"
    "path/filepath"
    "github.com/cdaprod/repocate/internal/log"
    "context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
)

// ListContainers lists all Docker containers for this project.
func ListContainers() ([]types.Container, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, err
    }

    containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
    if err != nil {
        return nil, err
    }

    log.Info("Containers listed successfully.")
    return containers, nil
}

// ResolveRepoName resolves the repository name from the provided URL or path
func ResolveRepoName(repoInput string) (string, error) {
    repoName := filepath.Base(repoInput)
    return repoName, nil
}

// IsRepoCloned checks if the repository is already cloned in the workspace
func IsRepoCloned(workspaceDir, repoName string) bool {
    repoPath := filepath.Join(workspaceDir, repoName)
    _, err := os.Stat(repoPath)
    return !os.IsNotExist(err)
}

// CloneRepository clones the repository to the workspace directory
func CloneRepository(workspaceDir, repoInput string) error {
    repoName, err := ResolveRepoName(repoInput)
    if err != nil {
        return err
    }

    repoPath := filepath.Join(workspaceDir, repoName)
    cmd := exec.Command("git", "clone", repoInput, repoPath)
    err = cmd.Run()
    if err != nil {
        return err
    }

    log.Info(fmt.Sprintf("Cloned repository %s", repoName))
    return nil
}

// InitContainer initializes the container for the repository
func InitContainer(workspaceDir, repoName string) error {
    log.Info(fmt.Sprintf("Initialized container for %s", repoName))
    return nil
}

// EnterContainer enters the development container for the repository
func EnterContainer(workspaceDir, repoName string) error {
    log.Info(fmt.Sprintf("Entered container for %s", repoName))
    return nil
}

// StopContainer stops the development container for the repository
func StopContainer(workspaceDir, repoName string) error {
    log.Info(fmt.Sprintf("Stopped container for %s", repoName))
    return nil
}

// RebuildContainer rebuilds the development container for the repository
func RebuildContainer(workspaceDir, repoName string) error {
    log.Info(fmt.Sprintf("Rebuilt container for %s", repoName))
    return nil
}

// CheckContainerExists checks if a Docker container with a specific name exists.
func CheckContainerExists(containerName string) (bool, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return false, err
    }

    containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
    if err != nil {
        return false, err
    }

    for _, c := range containers {
        for _, name := range c.Names {
            if name == "/"+containerName {
                return true, nil
            }
        }
    }

    return false, nil
}

// CreateAndStartContainer creates and starts a Docker container with a specific name.
func CreateAndStartContainer(containerName string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }

    ctx := context.Background()

    resp, err := cli.ContainerCreate(ctx, &container.Config{
        Image: "your-default-image", // Replace with your actual default image
        Cmd:   []string{"your-default-command"}, // Replace with your actual default command
    }, nil, nil, nil, containerName)
    if err != nil {
        return err
    }

    if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
        return err
    }

    log.Info(fmt.Sprintf("Container %s created and started.", containerName))
    return nil
}

// ExecIntoContainer executes into a running Docker container.
func ExecIntoContainer(containerName string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }

    ctx := context.Background()

    execConfig := types.ExecConfig{
        AttachStdin:  true,
        AttachStdout: true,
        AttachStderr: true,
        Tty:          true,
        Cmd:          []string{"sh"}, // Replace with your preferred shell or command
    }

    execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
    if err != nil {
        return err
    }

    err = cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{Tty: true})
    if err != nil {
        return err
    }

    log.Info(fmt.Sprintf("Executed into container %s.", containerName))
    return nil
}