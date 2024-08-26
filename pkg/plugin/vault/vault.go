package vault

import (
    "github.com/hashicorp/vault/api"
    "os"
)

// NewClient initializes and returns a new Vault client
func NewClient() (*api.Client, error) {
    config := api.DefaultConfig()
    if addr := os.Getenv("VAULT_ADDR"); addr != "" {
        config.Address = addr
    }

    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }

    return client, nil
}
