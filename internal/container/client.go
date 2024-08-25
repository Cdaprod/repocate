// internal/container/client.go
package container

import (
    "fmt"
    "github.com/cdaprod/repocate/internal/log"
    "github.com/docker/docker/client"
    "context"
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

// Add other utility functions related to Docker client here, if any...