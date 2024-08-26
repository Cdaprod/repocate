package tailscale

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
)

type ComposeService interface {
    AddService(serviceName string, serviceConfig interface{}) error
    RemoveService(serviceName string) error
    SaveToFile(fileName string) error
}

// AbstractDockerCompose represents a more abstracted Docker Compose structure
type AbstractDockerCompose struct {
    version  string
    services map[string]interface{}
}

type DockerCompose struct {
    Version  string                `yaml:"version"`
    Services map[string]ServiceDef `yaml:"services"`
}

type ServiceDef struct {
    Image       string            `yaml:"image,omitempty"`
    NetworkMode string            `yaml:"network_mode,omitempty"`
    Environment map[string]string `yaml:"environment,omitempty"`
    // Add more fields as needed
}

// NewAbstractDockerCompose initializes a new abstract Docker Compose
func NewAbstractDockerCompose(version string) *AbstractDockerCompose {
    return &AbstractDockerCompose{
        version:  version,
        services: make(map[string]interface{}),
    }
}

// AddService adds a new service to the abstract Docker Compose
func (adc *AbstractDockerCompose) AddService(serviceName string, serviceConfig interface{}) error {
    if serviceName == "" {
        return errors.New("service name cannot be empty")
    }
    adc.services[serviceName] = serviceConfig
    return nil
}

// SaveToFile saves the abstract Docker Compose to a file
func (adc *AbstractDockerCompose) SaveToFile(fileName string) error {
    data, err := yaml.Marshal(adc)
    if err != nil {
        return err
    }
    return ioutil.WriteFile(fileName, data, 0644)
}

// AddTailscaleSidecar adds the Tailscale sidecar to the Docker Compose file
func AddTailscaleSidecar(composeFile string, config *Config) error {
    data, err := ioutil.ReadFile(composeFile)
    if err != nil {
        return err
    }

    var compose DockerCompose
    err = yaml.Unmarshal(data, &compose)
    if err != nil {
        return err
    }

    // Define the Tailscale sidecar service
    tsService := ServiceDef{
        Image: "tailscale/tailscale:latest",
        Environment: map[string]string{
            "TS_AUTHKEY": config.AuthKey,
            "TS_EXIT_NODE": func() string {
                if config.ExitNode {
                    return "true"
                }
                return "false"
            }(),
        },
        NetworkMode: "service:" + config.ServiceName,
    }

    // Add the Tailscale service
    compose.Services[config.TailscaleContainer] = tsService

    // Marshal the updated compose structure back to YAML
    updatedData, err := yaml.Marshal(&compose)
    if err != nil {
        return err
    }

    // Write the updated YAML back to the file
    err = ioutil.WriteFile(composeFile, updatedData, 0644)
    if err != nil {
        return err
    }

    return nil
}
