package container

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cdaprod/repocate/internal/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
)

// handleDefaultContainer handles the default container initialization and startup
func handleDefaultContainer() {
    color.Cyan("Initializing and starting the 'repocate-default' container...")

    showProgress("Checking container status...", 100)

    // Initialize the default container if not exists or ensure it is running
    err := container.InitRepocateDefaultContainer()
    if err != nil {
        fmt.Println(color.RedString("Error initializing 'repocate-default' container: %s", err))
        os.Exit(1)
    }

    color.Green("Checking status of the 'repocate-default' container...")

    // Check if the container is running
    isRunning, err := container.IsContainerRunning("repocate-default")
    if err != nil {
        fmt.Println(color.RedString("Error checking container status: %s", err))
        os.Exit(1)
    }

    if !isRunning {
        color.Yellow("Container 'repocate-default' is not running. Starting it now...")

        err := container.StartContainer("repocate-default")
        if err != nil {
            fmt.Println(color.RedString("Error starting container: %s", err))
            os.Exit(1)
        }
    }

    color.Green("'repocate-default' container is ready.")
// initializeDockerClient initializes the Docker client and handles any errors.
func initializeDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create Docker client: %s", err))
		return nil, err
	}
	return cli, nil
}

// InitRepocateDefaultContainer initializes the 'repocate-default' container if it doesn't exist.
func InitRepocateDefaultContainer() error {
	containerName := GetDefaultContainerName()
	imageName := GetDefaultImageName()

	exists, err := CheckContainerExists(containerName)
	if err != nil {
		return fmt.Errorf("failed to check container existence: %w", err)
	}

	if exists {
		color.Green("Default container '%s' exists. Checking status...", containerName)

		isRunning, err := IsContainerRunning(containerName)
		if err != nil {
			return fmt.Errorf("failed to check if container is running: %w", err)
		}

		if !isRunning {
			color.Yellow("Container '%s' is not running. Starting it now...", containerName)
			if err := StartContainer(containerName); err != nil {
				return fmt.Errorf("failed to start container: %w", err)
			}
		}

		color.Green("Container '%s' is ready.", containerName)
		return nil
	}

	color.Yellow("Default container '%s' not found. Pulling image '%s'...", containerName, imageName)
	if err := PullImage(imageName); err != nil {
		return fmt.Errorf("failed to pull image '%s': %w", imageName, err)
	}

	color.Yellow("Creating and starting container '%s'...", containerName)
	if err := CreateAndStartContainer(containerName, imageName, []string{"/bin/zsh"}); err != nil {
		return fmt.Errorf("failed to create and start container '%s': %w", containerName, err)
	}

	color.Green("Default container '%s' created and started successfully.", containerName)
	return nil
}

// CheckContainerExists checks if a Docker container with a specific name exists.
func CheckContainerExists(containerName string) (bool, error) {
	cli, err := initializeDockerClient()
	if err != nil {
		return false, err
	}
	defer cli.Close()

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

// IsContainerRunning checks if a Docker container with a specific name is running.
func IsContainerRunning(containerName string) (bool, error) {
	cli, err := initializeDockerClient()
	if err != nil {
		return false, err
	}
	defer cli.Close()

	ctx := context.Background()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return false, err
	}

	for _, c := range containers {
		if c.Names[0] == "/"+containerName && c.State == "running" {
			return true, nil
		}
	}

	return false, nil
}

// StartContainer starts a Docker container with a specific name.
func StartContainer(containerName string) error {
	cli, err := initializeDockerClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx := context.Background()

	if err := cli.ContainerStart(ctx, containerName, types.ContainerStartOptions{}); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Container %s started successfully.", containerName))
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

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   cmd,
	}, nil, nil, nil, containerName)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create container %s: %s", containerName, err))
		return err
	}

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
		Cmd:          []string{"/bin/zsh"},
	}

	execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return fmt.Errorf("Failed to create exec configuration: %w", err)
	}

	resp, err := cli.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		return fmt.Errorf("Failed to attach to container exec process: %w", err)
	}
	defer resp.Close()

	_, err = io.Copy(os.Stdout, resp.Reader)
	if err != nil {
		return fmt.Errorf("Error during exec process copy: %w", err)
	}

	log.Info(fmt.Sprintf("Executed into container %s.", containerName))
	return nil
}

// StopContainer stops a running Docker container.
func StopContainer(containerName string) error {
	cli, err := initializeDockerClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx := context.Background()

	timeout := 10 // seconds
	if err := cli.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout}); err != nil {
		return fmt.Errorf("Failed to stop container %s: %w", containerName, err)
	}

	log.Info(fmt.Sprintf("Container %s stopped successfully.", containerName))
	return nil
}

// RemoveContainer removes a Docker container.
func RemoveContainer(containerName string) error {
	cli, err := initializeDockerClient()
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx := context.Background()

	if err := cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true}); err != nil {
		return fmt.Errorf("Failed to remove container %s: %w", containerName, err)
	}

	log.Info(fmt.Sprintf("Container %s removed successfully.", containerName))
	return nil
}

// ListContainers lists all Docker containers for this project.
func ListContainers() ([]types.Container, error) {
	cli, err := initializeDockerClient()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return containers, nil
}