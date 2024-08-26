package vault

import (
    "github.com/hashicorp/vault/api"
    "os"
)

// AuthenticateWithAppRole authenticates using the AppRole auth method
func AuthenticateWithAppRole(client *api.Client) error {
    roleID := os.Getenv("VAULT_ROLE_ID")
    secretID := os.Getenv("VAULT_SECRET_ID")

    data := map[string]interface{}{
        "role_id":   roleID,
        "secret_id": secretID,
    }

    secret, err := client.Logical().Write("auth/approle/login", data)
    if err != nil {
        return err
    }

    client.SetToken(secret.Auth.ClientToken)
    return nil
}
