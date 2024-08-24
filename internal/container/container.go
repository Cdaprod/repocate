// internal/container/container.go
package container

import (
    "os"
    "io"
    "fmt"
    "os/exec"
    "path/filepath"
    "github.com/cdaprod/repocate/internal/log"
    "context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/client"
)

// initializeDockerClient initializes the Docker client and handles any errors.
func initializeDockerClient() (*client.Client, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        log.Error(fmt.Sprintf("Failed to create Docker client: %s", err))
        return nil, err
    }
    return cli, nil
}

// ListContainers lists all Docker containers for this project.
func ListContainers() ([]types.Container, error) {
    cli, err := initializeDockerClient()
    if err != nil {
        return nil, err
    }
    defer cli.Close()

    containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
    if err != nil {
        return nil, err
    }

    log.Info("Containers listed successfully.")
    return containers, nil
}

// ResolveRepoName resolves the repository name from the provided URL or path.
func ResolveRepoName(repoInput string) (string, error) {
    repoName := filepath.Base(repoInput)
    return repoName, nil
}

// IsRepoCloned checks if the repository is already cloned in the workspace.
func IsRepoCloned(workspaceDir, repoName string) bool {
    repoPath := filepath.Join(workspaceDir, repoName)
    _, err := os.Stat(repoPath)
    return !os.IsNotExist(err)
}

// CloneRepository clones the repository to the workspace directory.
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

// InitContainer initializes the container for the repository.
func InitContainer(workspaceDir, repoName string) error {
    exists, err := CheckContainerExists(repoName)
    if err != nil {
        return err
    }

    if exists {
        log.Info(fmt.Sprintf("Container for %s already exists. Entering container.", repoName))
        return ExecIntoContainer(repoName)
    }

    log.Info(fmt.Sprintf("Initializing container for %s", repoName))
    err = CreateAndStartContainer(repoName, "cdaprod/repocate-dev:v1.0.0-arm64", []string{"/bin/zsh"})
    if err != nil {
        return err
    }

    return ExecIntoContainer(repoName)
}

// EnterContainer enters the development container for the repository.
func EnterContainer(workspaceDir, repoName string) error {
    log.Info(fmt.Sprintf("Entering container for %s", repoName))
    return ExecIntoContainer(repoName)
}

// StopContainer stops the development container for the repository
func StopContainer(workspaceDir, repoName string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }

    ctx := context.Background()

    // StopOptions requires a Timeout
    stopOptions := container.StopOptions{
        Timeout: nil, // or you can specify a time.Duration value like `new(int)` with a value for timeout in seconds
    }

    err = cli.ContainerStop(ctx, repoName, stopOptions)
    if err != nil {
        log.Error(fmt.Sprintf("Failed to stop container %s: %s", repoName, err))
        return err
    }

    log.Info(fmt.Sprintf("Stopped container for %s", repoName))
    return nil
}

// RebuildContainer rebuilds the development container for the repository.
func RebuildContainer(workspaceDir, repoName string) error {
    log.Info(fmt.Sprintf("Rebuilding container for %s", repoName))
    err := StopContainer(workspaceDir, repoName)
    if err != nil {
        return err
    }

    err = InitContainer(workspaceDir, repoName)
    if err != nil {
        return err
    }

    log.Info(fmt.Sprintf("Rebuilt container for %s", repoName))
    return nil
}

// CheckImageExists checks if a Docker image with a specific name exists locally.
func CheckImageExists(imageName string) (bool, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return false, err
    }
    defer cli.Close()

    ctx := context.Background()

    _, _, err = cli.ImageInspectWithRaw(ctx, imageName)
    if err == nil {
        // Image already exists locally
        log.Info(fmt.Sprintf("Image %s already exists locally.", imageName))
        return true, nil
    }

    if client.IsErrNotFound(err) {
        log.Info(fmt.Sprintf("Image %s not found locally.", imageName))
        return false, nil
    }

    // Handle any other error
    return false, err
}

// PullImage pulls a Docker image, ensuring it is present locally.
func PullImage(imageName string) error {
    cli, err := initializeDockerClient()
    if err != nil {
        return fmt.Errorf("failed to create Docker client: %w", err)
    }
    defer cli.Close()

    ctx := context.Background()

    // Check if the image exists locally
    exists, err := CheckImageExists(imageName)
    if err != nil {
        return fmt.Errorf("error checking image existence: %w", err)
    }

    if exists {
        log.Info(fmt.Sprintf("Image %s already exists locally. Skipping pull.", imageName))
        return nil
    }

    // If image does not exist, pull it
    out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
    if err != nil {
        return fmt.Errorf("failed to pull image %s: %w", imageName, err)
    }
    defer out.Close()

    log.Info(fmt.Sprintf("Pulling image %s...", imageName))
    if _, err := io.Copy(os.Stdout, out); err != nil {
        return fmt.Errorf("failed to read image pull response: %w", err)
    }

    return nil
}

// CreateAndStartContainer creates and starts a Docker container with a specific name.
func CreateAndStartContainer(containerName, imageName string, cmd []string) error {
    cli, err := initializeDockerClient()
    if err != nil {
        return err
    }
    defer cli.Close()

    ctx := context.Background()

    // Ensure the image is pulled or exists locally
    err = PullImage(imageName)
    if err != nil {
        log.Error(fmt.Sprintf("Error pulling image: %s", err))
        return err
    }

    // Create Docker container
    resp, err := cli.ContainerCreate(ctx, &container.Config{
        Image: imageName,
        Cmd:   cmd,
    }, nil, nil, nil, containerName)
    if err != nil {
        log.Error(fmt.Sprintf("Failed to create container %s: %s", containerName, err))
        return err
    }

    // Start Docker container
    if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
        log.Error(fmt.Sprintf("Failed to start container %s: %s", containerName, err))
        return err
    }

    log.Info(fmt.Sprintf("Container %s created and started successfully with image %s.", containerName, imageName))
    return nil
}

// ExecIntoContainer executes into a running Docker container.
func ExecIntoContainer(containerName string) error {
    cli, err := initializeDockerClient()
    if err != nil {
        return err
    }
    defer cli.Close()

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