package tailscale

// GenericServiceConfig represents a generic configuration for any service
type GenericServiceConfig[T any] struct {
    ServiceName   string
    ConfigDetails T
}

type Config struct {
    AuthKey            string
    ExitNode           bool
    Funnel             bool
    ServeEnabled       bool
    ServiceName        string
    TailscaleContainer string
}

// NewConfig creates a new Tailscale configuration
func NewConfig(authKey string, exitNode, funnel, serveEnabled bool, serviceName, containerName string) *Config {
    return &Config{
        AuthKey:            authKey,
        ExitNode:           exitNode,
        Funnel:             funnel,
        ServeEnabled:       serveEnabled,
        ServiceName:        serviceName,
        TailscaleContainer: containerName,
    }
}

// NewGenericServiceConfig creates a new generic service configuration
func NewGenericServiceConfig[T any](serviceName string, configDetails T) *GenericServiceConfig[T] {
    return &GenericServiceConfig[T]{
        ServiceName:   serviceName,
        ConfigDetails: configDetails,
    }
}
