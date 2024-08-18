package container

import (
    "context"
    "fmt"
    "os"
    "os/exec"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/api/types/mount"
    "github.com/docker/docker/client"
    "repocate/internal/log"
    "repocate/internal/utils"
)

// InitContainer initializes a Docker container for the repo.
func InitContainer(workspaceDir, repoName string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }

    containerName := fmt.Sprintf("repocate-%s", repoName)
    hostPath := utils.GetRepoPath(workspaceDir, repoName)
    containerPath := "/workspace"

    _, err = cli.ContainerCreate(context.Background(), &container.Config{
        Image: "repocate-base-image",
        Cmd:   []string{"tail", "-f", "/dev/null"},
    }, &container.HostConfig{
        Mounts: []mount.Mount{
            {
                Type:   mount.TypeBind,
                Source: hostPath,
                Target: containerPath,
            },
        },
    }, nil, nil, containerName)
    if err != nil {
        return err
    }

    log.Info("Container initialized successfully")
    return nil
}

// EnterContainer allows the user to enter an existing container.
func EnterContainer(workspaceDir, repoName string) error {
    containerName := fmt.Sprintf("repocate-%s", repoName)

    cmd := exec.Command("docker", "exec", "-it", containerName, "/bin/bash")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    return cmd.Run()
}

// StopContainer stops a running container.
func StopContainer(repoName string) error {
    containerName := fmt.Sprintf("repocate-%s", repoName)
    cmd := exec.Command("docker", "stop", containerName)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Run()
}

// RebuildContainer rebuilds a Docker container.
func RebuildContainer(workspaceDir, repoName string) error {
    if err := StopContainer(repoName); err != nil {
        return err
    }
    cmd := exec.Command("docker", "rm", fmt.Sprintf("repocate-%s", repoName))
    if err := cmd.Run(); err != nil {
        return err
    }
    return InitContainer(workspaceDir, repoName)
}