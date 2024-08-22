// internal/container/volumes.go
package container

import (
    "context"
//    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/filters"
    "github.com/docker/docker/api/types/volume"
    "github.com/docker/docker/client"
    "github.com/cdaprod/repocate/internal/log"
)

// CreateVolume creates a Docker volume for persistent storage.
func CreateVolume(volumeName string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }

    _, err = cli.VolumeCreate(context.Background(), volume.VolumeCreateBody{
        Name: volumeName,
    })
    if err != nil {
        return err
    }

    log.Info("Volume created successfully")
    return nil
}

// RemoveVolume removes a Docker volume.
func RemoveVolume(volumeName string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }

    err = cli.VolumeRemove(context.Background(), volumeName, true)
    if err != nil {
        return err
    }

    log.Info("Volume removed successfully")
    return nil
}

// ListVolumes lists all Docker volumes associated with Repocate.
// ListVolumes lists all Docker volumes associated with Repocate.
func ListVolumes() ([]volume.Volume, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, err
    }

    volumeList, err := cli.VolumeList(context.Background(), filters.Args{})
    if err != nil {
        return nil, err
    }

    // Convert []*volume.Volume to []volume.Volume
    var volumes []volume.Volume
    for _, v := range volumeList.Volumes {
        volumes = append(volumes, *v)
    }

    return volumes, nil
}