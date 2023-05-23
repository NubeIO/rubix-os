package compat

import (
	"github.com/NubeIO/rubix-os/plugin/pluginapi"
	"github.com/gin-gonic/gin"
	"net/url"
)

// Capability is a capability the plugin provides.
type Capability string

const (
	// Configurer are consigurables.
	Configurer = Capability("configurer")
	// Storager stores data.
	Storager = Capability("storager")
	// Webhooker registers webhooks.
	Webhooker = Capability("webhooker")
	// Displayer displays instructions.
	Displayer = Capability("displayer")
)

// PluginInstance is an encapsulation layer of plugin instances of different backends.
type PluginInstance interface {
	Enable() error
	Disable() error
	// GetDisplay see Displayer
	GetDisplay(location *url.URL) pluginapi.Response
	// DefaultConfig see Configurer
	DefaultConfig() interface{}
	// GetConfig see Configurer
	GetConfig() interface{}
	// ValidateAndSetConfig see Configurer
	ValidateAndSetConfig(c interface{}) error

	// RegisterWebhook see Webhooker#RegisterWebhook
	RegisterWebhook(basePath string, mux *gin.RouterGroup)

	// SetStorageHandler see Storager#SetStorageHandler.
	SetStorageHandler(handler StorageHandler)

	// Supports Returns the supported modules, f.ex. storager
	Supports() Capabilities
}

// HasSupport tests a PluginInstance for a capability.
func HasSupport(p PluginInstance, toCheck Capability) bool {
	for _, module := range p.Supports() {
		if module == toCheck {
			return true
		}
	}
	return false
}

// Capabilities is a slice of module.
type Capabilities []Capability

// Strings converts []Module to []string.
func (m Capabilities) Strings() []string {
	var result []string
	for _, module := range m {
		result = append(result, string(module))
	}
	return result
}

// StorageHandler see plugin.StorageHandler.
type StorageHandler interface {
	Save(b []byte) error
	Load() ([]byte, error)
}
