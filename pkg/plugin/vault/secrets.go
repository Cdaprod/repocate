package vault

import (
    "github.com/hashicorp/vault/api"
)

// GetSecret retrieves a secret from Vault at the specified path
func GetSecret(client *api.Client, path string) (map[string]interface{}, error) {
    secret, err := client.Logical().Read(path)
    if err != nil {
        return nil, err
    }

    if secret == nil || secret.Data == nil {
        return nil, fmt.Errorf("no data found at path: %s", path)
    }

    return secret.Data, nil
}
