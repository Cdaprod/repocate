package tailscale

import (
    "errors"
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

// DockerCompose represents a Docker Compose configuration
type DockerCompose struct {
    Version  string               `yaml:"version"`
    Services map[string]ServiceDef `yaml:"services"`
}

// ServiceDef defines a service in Docker Compose
type ServiceDef struct {
    Image       string            `yaml:"image,omitempty"`
    NetworkMode string            `yaml:"network_mode,omitempty"`
    Environment map[string]string `yaml:"environment,omitempty"`
}

// AddTailscaleSidecar adds a Tailscale sidecar to a Docker Compose configuration
func AddTailscaleSidecar(composeFile string, authKey string, serviceName string) error {
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
            "TS_AUTHKEY": authKey,
            "TS_EXIT_NODE": "false",
        },
        NetworkMode: "service:" + serviceName,
    }

    // Add the Tailscale service to the Docker Compose
    compose.Services["tailscale-sidecar"] = tsService

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
