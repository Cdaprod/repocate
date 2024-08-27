package registry

import (
	"fmt"
	"sync"

	"github.com/cdaprod/repocate/internal/log"
)

// PluginHandler is a function type for plugin handlers
type PluginHandler func() error

// ImageConfig holds configuration for a Docker image
type ImageConfig struct {
	Name    string
	Tag     string
	BuildArgs map[string]string
}

// ContainerConfig holds configuration for a container
type ContainerConfig struct {
	Image       string
	Entrypoint  []string
	Cmd         []string
	Environment map[string]string
	Volumes     []string
	Ports       []string
}

// Registry holds the registered components like plugins, services, etc.
type Registry struct {
	mu         sync.RWMutex
	plugins    map[string]PluginHandler
	images     map[string]ImageConfig
	containers map[string]ContainerConfig
}

// NewRegistry initializes a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins:    make(map[string]PluginHandler),
		images:     make(map[string]ImageConfig),
		containers: make(map[string]ContainerConfig),
	}
}

// RegisterPlugin adds a new plugin to the registry.
func (r *Registry) RegisterPlugin(name string, handler PluginHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.plugins[name]; exists {
		return fmt.Errorf("plugin already registered: %s", name)
	}
	r.plugins[name] = handler
	log.Info(fmt.Sprintf("Registered plugin: %s", name))
	return nil
}

// RegisterImage adds a new Docker image to the registry.
func (r *Registry) RegisterImage(config ImageConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.images[config.Name]; exists {
		return fmt.Errorf("image already registered: %s", config.Name)
	}
	r.images[config.Name] = config
	log.Info(fmt.Sprintf("Registered image: %s with tag: %s", config.Name, config.Tag))
	return nil
}

// RegisterContainer adds a new container to the registry.
func (r *Registry) RegisterContainer(name string, config ContainerConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.containers[name]; exists {
		return fmt.Errorf("container already registered: %s", name)
	}
	r.containers[name] = config
	log.Info(fmt.Sprintf("Registered container: %s with image: %s", name, config.Image))
	return nil
}

// GetPlugin retrieves a plugin handler by name.
func (r *Registry) GetPlugin(name string) (PluginHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, exists := r.plugins[name]
	return handler, exists
}

// GetImage retrieves an image configuration by name.
func (r *Registry) GetImage(name string) (ImageConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	image, exists := r.images[name]
	return image, exists
}

// GetContainer retrieves a container configuration by name.
func (r *Registry) GetContainer(name string) (ContainerConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	container, exists := r.containers[name]
	return container, exists
}

// UnregisterPlugin removes a plugin from the registry.
func (r *Registry) UnregisterPlugin(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.plugins[name]; exists {
		delete(r.plugins, name)
		log.Info(fmt.Sprintf("Unregistered plugin: %s", name))
	}
}

// UnregisterImage removes an image from the registry.
func (r *Registry) UnregisterImage(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.images[name]; exists {
		delete(r.images, name)
		log.Info(fmt.Sprintf("Unregistered image: %s", name))
	}
}

// UnregisterContainer removes a container from the registry.
func (r *Registry) UnregisterContainer(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.containers[name]; exists {
		delete(r.containers, name)
		log.Info(fmt.Sprintf("Unregistered container: %s", name))
	}
}

// ListPlugins returns a list of all registered plugin names.
func (r *Registry) ListPlugins() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	plugins := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		plugins = append(plugins, name)
	}
	return plugins
}

// ListImages returns a list of all registered image names.
func (r *Registry) ListImages() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	images := make([]string, 0, len(r.images))
	for name := range r.images {
		images = append(images, name)
	}
	return images
}

// ListContainers returns a list of all registered container names.
func (r *Registry) ListContainers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	containers := make([]string, 0, len(r.containers))
	for name := range r.containers {
		containers = append(containers, name)
	}
	return containers
}

// ExecutePlugin executes a plugin by name.
func (r *Registry) ExecutePlugin(name string) error {
	handler, exists := r.GetPlugin(name)
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}
	return handler()
}