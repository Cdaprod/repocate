package container

import (
    "fmt"
    "github.com/docker/docker/api/types"
)

// ContainerConfig represents the configuration for a container
type ContainerConfig struct {
    Name  string
    Image string
    Cmd   []string
    // Add other fields as needed
}

// ContainerInfo represents information about a container
type ContainerInfo struct {
    ID     string
    Name   string
    Status string
    // Add other fields as needed
}

// DockerClientInterface defines the methods that a Docker client should implement
type DockerClientInterface interface {
    CreateContainer(config ContainerConfig) (string, error)
    StartContainer(id string) error
    StopContainer(id string) error
    PullImage(name string) error          // Added for image operations
    CheckImageExists(name string) (bool, error) // Added for image operations
    ListImages() ([]types.ImageSummary, error)  // Added for image operations
    RemoveImage(name string) error        // Added for image operations
    // Add other methods as needed
}

// Some constants
const (
    DefaultImageName     = "cdaprod/repocate-dev:1.0.0-arm64"
    DefaultContainerName = "repocate-default"
)

// Custom error types
type ErrContainerNotFound struct {
    Name string
}

func (e ErrContainerNotFound) Error() string {
    return fmt.Sprintf("container not found: %s", e.Name)
}

// Custom error types for image operations
type ErrImageNotFound struct {
    Name string
}

func (e ErrImageNotFound) Error() string {
    return fmt.Sprintf("image not found: %s", e.Name)
}