package plugin

import (
	"fmt"
	"sync"
)

// Plugin interface that all plugins must implement
type Plugin interface {
	Execute() error
	Name() string
}

// Registry to hold registered plugins
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]Plugin
}

// NewRegistry initializes a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
	}
}

// RegisterPlugin registers a new plugin in the registry
func (r *Registry) RegisterPlugin(p Plugin) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[p.Name()] = p
	fmt.Printf("Registered plugin: %s\n", p.Name())
}

// InjectDependencies allows injecting dependencies into a plugin
func (r *Registry) InjectDependencies(name string, deps ...interface{}) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}

	// Perform type assertions and inject dependencies based on plugin type
	switch p := plugin.(type) {
	case *TailscalePlugin:
		if len(deps) >= 2 {
			p.authKey, _ = deps[0].(string)
			p.serviceName, _ = deps[1].(string)
		}
	case *VaultPlugin:
		if len(deps) >= 1 {
			p.url, _ = deps[0].(string)
		}
		// Add more cases for other plugins as needed
	default:
		return fmt.Errorf("unknown plugin type: %T", p)
	}

	return nil
}

// GetPlugin retrieves a plugin by name
func (r *Registry) GetPlugin(name string) (Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	plugin, exists := r.plugins[name]
	return plugin, exists
}

// ExecutePlugin executes a plugin by name
func (r *Registry) ExecutePlugin(name string) error {
	plugin, exists := r.GetPlugin(name)
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}
	return plugin.Execute()
}