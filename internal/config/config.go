package config

import (
    "os"
    "path/filepath"
    "encoding/json"
    "fmt"
    "github.com/cdaprod/repocate/internal/log"
)

var (
    ConfigFile   = "repocate.json"  // Config file name
    WorkspaceDir = ""               // Path to workspace directory
)

// Config represents the structure of the configuration file
type Config struct {
    WorkspaceDir string `json:"workspace_dir"`
}

// LoadConfig loads configuration from the config file
func LoadConfig() {
    configPath := filepath.Join(getConfigDir(), ConfigFile)
    file, err := os.Open(configPath)
    if err != nil {
        log.Error(fmt.Sprintf("Could not open config file: %s", err))
        os.Exit(1)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    config := Config{}
    if err := decoder.Decode(&config); err != nil {
        log.Error(fmt.Sprintf("Error decoding config file: %s", err))
        os.Exit(1)
    }

    WorkspaceDir = config.WorkspaceDir
}

// SaveConfig saves the current configuration to the config file
func SaveConfig() {
    configPath := filepath.Join(getConfigDir(), ConfigFile)
    file, err := os.Create(configPath)
    if err != nil {
        log.Error(fmt.Sprintf("Could not create config file: %s", err))
        os.Exit(1)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    config := Config{WorkspaceDir: WorkspaceDir}
    if err := encoder.Encode(&config); err != nil {
        log.Error(fmt.Sprintf("Error encoding config file: %s", err))
        os.Exit(1)
    }
}

// getConfigDir returns the path to the configuration directory
func getConfigDir() string {
    configDir := filepath.Join(os.Getenv("HOME"), ".config", "repocate")
    if _, err := os.Stat(configDir); os.IsNotExist(err) {
        os.MkdirAll(configDir, 0755)
    }
    return configDir
}