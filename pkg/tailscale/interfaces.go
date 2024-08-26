package tailscale

// DockerService defines behavior for a Docker service
type DockerService interface {
    AddService(name string, config interface{}) error
    UpdateService(name string, config interface{}) error
    RemoveService(name string) error
    ListServices() ([]string, error)
}

// TailscaleManager defines behavior for managing Tailscale integration
type TailscaleManager interface {
    StartContainer(containerName string) error
    StopContainer(containerName string) error
    ConfigureService(serviceConfig interface{}) error
}
