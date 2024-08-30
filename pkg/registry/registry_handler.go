package registry

import (
	"fmt"
	"github.com/Cdaprod/registry-service/internal/registry"
	"github.com/Cdaprod/registry-service/pkg/plugins"
)

type RegistryHandler struct {
	itemStore    *registry.ItemStore
	pluginLoader *plugins.PluginLoader
}

func NewRegistryHandler(pluginsDir string) *RegistryHandler {
	itemStore := registry.NewItemStore()
	return &RegistryHandler{
		itemStore:    itemStore,
		pluginLoader: plugins.NewPluginLoader(itemStore, pluginsDir),
	}
}

func (rh *RegistryHandler) RegisterPlugin(name string, handler PluginHandler) error {
	item := &registry.Item{
		Type: "Plugin",
		Name: name,
		Metadata: map[string]interface{}{
			"handler": handler,
		},
	}
	_, err := rh.itemStore.UpsertItem(item)
	return err
}

func (rh *RegistryHandler) RegisterImage(config ImageConfig) error {
	item := &registry.Item{
		Type: "Image",
		Name: config.Name,
		Metadata: map[string]interface{}{
			"tag":       config.Tag,
			"build_args": config.BuildArgs,
		},
	}
	_, err := rh.itemStore.UpsertItem(item)
	return err
}

func (rh *RegistryHandler) RegisterContainer(name string, config ContainerConfig) error {
	item := &registry.Item{
		Type: "Container",
		Name: name,
		Metadata: map[string]interface{}{
			"image":       config.Image,
			"entrypoint":  config.Entrypoint,
			"cmd":         config.Cmd,
			"environment": config.Environment,
			"volumes":     config.Volumes,
			"ports":       config.Ports,
		},
	}
	_, err := rh.itemStore.UpsertItem(item)
	return err
}

func (rh *RegistryHandler) GetPlugin(name string) (PluginHandler, bool) {
	item, err := rh.itemStore.GetItem(name)
	if err != nil {
		return nil, false
	}
	if item.Type != "Plugin" {
		return nil, false
	}
	handler, ok := item.Metadata["handler"].(PluginHandler)
	return handler, ok
}

func (rh *RegistryHandler) GetImage(name string) (ImageConfig, bool) {
	item, err := rh.itemStore.GetItem(name)
	if err != nil {
		return ImageConfig{}, false
	}
	if item.Type != "Image" {
		return ImageConfig{}, false
	}
	return ImageConfig{
		Name:      item.Name,
		Tag:       item.Metadata["tag"].(string),
		BuildArgs: item.Metadata["build_args"].(map[string]string),
	}, true
}

func (rh *RegistryHandler) GetContainer(name string) (ContainerConfig, bool) {
	item, err := rh.itemStore.GetItem(name)
	if err != nil {
		return ContainerConfig{}, false
	}
	if item.Type != "Container" {
		return ContainerConfig{}, false
	}
	return ContainerConfig{
		Image:       item.Metadata["image"].(string),
		Entrypoint:  item.Metadata["entrypoint"].([]string),
		Cmd:         item.Metadata["cmd"].([]string),
		Environment: item.Metadata["environment"].(map[string]string),
		Volumes:     item.Metadata["volumes"].([]string),
		Ports:       item.Metadata["ports"].([]string),
	}, true
}

func (rh *RegistryHandler) UnregisterPlugin(name string) {
	rh.itemStore.DeleteItem(name)
}

func (rh *RegistryHandler) UnregisterImage(name string) {
	rh.itemStore.DeleteItem(name)
}

func (rh *RegistryHandler) UnregisterContainer(name string) {
	rh.itemStore.DeleteItem(name)
}

func (rh *RegistryHandler) ListPlugins() []string {
	items := rh.itemStore.ListItems()
	var plugins []string
	for _, item := range items {
		if item.Type == "Plugin" {
			plugins = append(plugins, item.Name)
		}
	}
	return plugins
}

func (rh *RegistryHandler) ListImages() []string {
	items := rh.itemStore.ListItems()
	var images []string
	for _, item := range items {
		if item.Type == "Image" {
			images = append(images, item.Name)
		}
	}
	return images
}

func (rh *RegistryHandler) ListContainers() []string {
	items := rh.itemStore.ListItems()
	var containers []string
	for _, item := range items {
		if item.Type == "Container" {
			containers = append(containers, item.Name)
		}
	}
	return containers
}

func (rh *RegistryHandler) ExecutePlugin(name string) error {
	handler, exists := rh.GetPlugin(name)
	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}
	return handler()
}

func (rh *RegistryHandler) LoadAllPlugins() error {
	return rh.pluginLoader.LoadAll()
}

func (rh *RegistryHandler) LoadPlugin(path string) error {
	return rh.pluginLoader.LoadPlugin(path)
}