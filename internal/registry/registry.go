package registry

import (
	"fmt"
	"sync"
)

// Registry holds the registered components like plugins, services, etc.
type Registry struct {
	mu       sync.RWMutex
	plugins  map[string]func() // Map of plugin names to handler functions
	images   map[string]string // Map of image names to Docker image strings
	containers map[string]string // Map of container names to container settings
}

// NewRegistry initializes a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins:  make(map[string]func()),
		images:   make(map[string]string),
		containers: make(map[string]string),
	}
}

// RegisterPlugin adds a new plugin to the registry.
func (r *Registry) RegisterPlugin(name string, handler func()) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.plugins[name] = handler
	fmt.Printf("Registered plugin: %s\n", name)
}

// RegisterImage adds a new Docker image to the registry.
func (r *Registry) RegisterImage(name, image string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.images[name] = image
	fmt.Printf("Registered image: %s with Docker image: %s\n", name, image)
}

// RegisterContainer adds a new container to the registry.
func (r *Registry) RegisterContainer(name, container string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.containers[name] = container
	fmt.Printf("Registered container: %s with settings: %s\n", name, container)
}

// GetPlugin retrieves a plugin handler by name.
func (r *Registry) GetPlugin(name string) (func(), bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, exists := r.plugins[name]
	return handler, exists
}

// GetImage retrieves an image by name.
func (r *Registry) GetImage(name string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	image, exists := r.images[name]
	return image, exists
}

// GetContainer retrieves a container by name.
func (r *Registry) GetContainer(name string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	container, exists := r.containers[name]
	return container, exists
}

// UnregisterPlugin removes a plugin from the registry.
func (r *Registry) UnregisterPlugin(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.plugins, name)
	fmt.Printf("Unregistered plugin: %s\n", name)
}

// UnregisterImage removes an image from the registry.
func (r *Registry) UnregisterImage(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.images, name)
	fmt.Printf("Unregistered image: %s\n", name)
}

// UnregisterContainer removes a container from the registry.
func (r *Registry) UnregisterContainer(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.containers, name)
	fmt.Printf("Unregistered container: %s\n", name)
}
