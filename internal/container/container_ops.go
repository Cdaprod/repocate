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

// InitRepocateDefaultContainer initializes the 'repocate-default' container if it doesn't exist.
func InitRepocateDefaultContainer() error {
	containerName := GetDefaultContainerName()
	imageName := GetDefaultImageName()

	// Check if the container exists
	exists, err := CheckContainerExists(containerName)
	if err != nil {
		return fmt.Errorf("failed to check container existence: %w", err)
	}

	if exists {
		color.Green("Default container '%s' exists. Checking status...", containerName)

		// Ensure the container is running
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

	// If the container doesn't exist, pull the image
	color.Yellow("Default container '%s' not found. Pulling image '%s'...", containerName, imageName)
	if err := PullImage(imageName); err != nil {
		return fmt.Errorf("failed to pull image '%s': %w", imageName, err)
	}

	// Create and start the container
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

	// Get container details
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

	// Start the container
	if err := cli.ContainerStart(ctx, containerName, types.ContainerStartOptions{}); err != nil {
		return err
	}

	log.Info(fmt.Sprintf("Container %s started successfully.", containerName))
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
		Cmd:          []string{"/bin/zsh"}, // Make sure this shell exists in the container
	}

	execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return fmt.Errorf("Failed to create exec configuration: %w", err)
	}

	// Start the exec process
	execStartCheck := types.ExecStartCheck{
		Tty: true,
	}

	resp, err := cli.ContainerExecAttach(ctx, execID.ID, execStartCheck)
	if err != nil {
		return fmt.Errorf("Failed to attach to container exec process: %w", err)
	}
	defer resp.Close()

	// Copy output to stdout and stderr
	_, err = io.Copy(os.Stdout, resp.Reader)
	if err != nil {
		return fmt.Errorf("Error during exec process copy: %w", err)
	}

	log.Info(fmt.Sprintf("Executed into container %s.", containerName))
	return nil
}