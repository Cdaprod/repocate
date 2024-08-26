package tailscale

import (
    "log"
    "os/exec"
)

// StartTailscaleContainer starts the Tailscale container using Docker
func StartTailscaleContainer(containerName string) error {
    cmd := exec.Command("docker", "start", containerName)
    if err := cmd.Run(); err != nil {
        log.Printf("Failed to start Tailscale container: %v", err)
        return err
    }
    log.Printf("Tailscale container %s started successfully", containerName)
    return nil
}

// ConfigureTailscaleService sets up Tailscale for the service
func ConfigureTailscaleService(composeFile string, config *Config) error {
    // Modify Docker Compose to add Tailscale
    if err := AddTailscaleSidecar(composeFile, config); err != nil {
        return err
    }

    // Start Tailscale sidecar
    return StartTailscaleContainer(config.TailscaleContainer)
}
