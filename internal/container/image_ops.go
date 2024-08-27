package container

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cdaprod/repocate/internal/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CheckImageExists checks if a Docker image with a specific name exists locally.
func CheckImageExists(imageName string) (bool, error) {
	cli, err := initializeDockerClient()
	if err != nil {
		return false, err
	}
	defer cli.Close()

	ctx := context.Background()

	_, _, err = cli.ImageInspectWithRaw(ctx, imageName)
	if err == nil {
		log.Info(fmt.Sprintf("Image %s already exists locally.", imageName))
		return true, nil
	}

	if client.IsErrNotFound(err) {
		log.Info(fmt.Sprintf("Image %s not found locally.", imageName))
		return false, nil
	}

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

	exists, err := CheckImageExists(imageName)
	if err != nil {
		return fmt.Errorf("error checking image existence: %w", err)
	}

	if exists {
		log.Info(fmt.Sprintf("Image %s already exists locally. Skipping pull.", imageName))
		return nil
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", imageName, err)
	}
	defer out.Close()

	log.Info(fmt.Sprintf("Pulling image %s...", imageName))
	_, err = io.Copy(os.Stdout, out)
	if err != nil {
		return fmt.Errorf("error copying image pull output to stdout: %w", err)
	}

	log.Info(fmt.Sprintf("Image %s pulled successfully.", imageName))
	return nil
}

// RemoveImage removes a Docker image from the local repository.
func RemoveImage(imageName string) error {
	cli, err := initializeDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()

	_, err = cli.ImageRemove(ctx, imageName, types.ImageRemoveOptions{
		Force:         false,
		PruneChildren: true,
	})
	if err != nil {
		return fmt.Errorf("failed to remove image %s: %w", imageName, err)
	}

	log.Info(fmt.Sprintf("Image %s removed successfully.", imageName))
	return nil
}

// ListImages lists all Docker images available locally.
func ListImages() ([]types.ImageSummary, error) {
	cli, err := initializeDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list Docker images: %w", err)
	}

	for _, image := range images {
		log.Info(fmt.Sprintf("Found image: %s", image.ID))
	}

	return images, nil
}

// BuildImage builds a Docker image from a Dockerfile.
func BuildImage(dockerfilePath, imageName string) error {
	cli, err := initializeDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()

	buildContext, err := os.Open(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to open Dockerfile: %w", err)
	}
	defer buildContext.Close()

	options := types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageName},
		Remove:     true,
	}

	resp, err := cli.ImageBuild(ctx, buildContext, options)
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying image build output to stdout: %w", err)
	}

	log.Info(fmt.Sprintf("Image %s built successfully.", imageName))
	return nil
}