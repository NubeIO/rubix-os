package plugin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"plugin"
	"strings"
	"sync"
	"time"

	"github.com/NubeIO/flow-framework/model"
	"github.com/NubeIO/flow-framework/plugin/compat"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// The Database interface for encapsulating database access.
type Database interface {
	GetPlugin(id string) (*model.PluginConf, error)
	GetPluginByPath(path string) (*model.PluginConf, error)
	CreatePlugin(p *model.PluginConf) error
	CreateMessage(message *model.Message) error
	UpdatePluginConf(p *model.PluginConf) error
}

// Notifier notifies when a new message was created.
type Notifier interface {
	Notify(message *model.MessageExternal)
}

// Manager is an encapsulating layer for plugins and manages all plugins and its instances.
type Manager struct {
	mutex     *sync.RWMutex
	instances map[string]compat.PluginInstance
	plugins   map[string]compat.Plugin
	messages  chan model.MessageExternal
	db        Database
	mux       *gin.RouterGroup
}

// NewManager created a Manager from configurations.
func NewManager(db Database, directory string, mux *gin.RouterGroup, notifier Notifier) (*Manager, error) {
	manager := &Manager{
		mutex:     &sync.RWMutex{},
		instances: map[string]compat.PluginInstance{},
		plugins:   map[string]compat.Plugin{},
		messages:  make(chan model.MessageExternal),
		db:        db,
		mux:       mux,
	}
	go func() {
		for {
			message := <-manager.messages
			internalMsg := &model.Message{
				Title:    message.Title,
				Priority: message.Priority,
				Date:     message.Date,
				Message:  message.Message,
			}

			if message.Extras != nil {
				internalMsg.Extras, _ = json.Marshal(message.Extras)
			}
			err := db.CreateMessage(internalMsg)
			if err != nil {
				log.Println("error on send message", internalMsg.Title)
			}
			message.ID = internalMsg.ID
			notifier.Notify(&message)
		}
	}()

	if err := manager.loadPlugins(directory); err != nil {
		return nil, err
	}
	if err := manager.initializePlugins(); err != nil {
		return nil, err
	}
	return manager, nil
}

// ErrAlreadyEnabledOrDisabled is returned on SetPluginEnabled call when a plugin is already enabled or disabled.
var ErrAlreadyEnabledOrDisabled = errors.New("config is already on your state")

// SetPluginEnabled sets the plugins enabled state.
func (m *Manager) SetPluginEnabled(pluginID string, enabled bool) error {
	instance, err := m.Instance(pluginID)
	if err != nil {
		return errors.New("instance not found")
	}
	conf, err := m.db.GetPlugin(pluginID)
	if err != nil {
		return err
	}
	if conf.Enabled == enabled {
		return ErrAlreadyEnabledOrDisabled
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if enabled {
		err = instance.Enable()
	} else {
		err = instance.Disable()
	}
	if err != nil {
		return err
	}
	if newConf, err := m.db.GetPlugin(pluginID); /* conf might be updated by instance */ err == nil {
		conf = newConf
	}
	conf.Enabled = enabled
	return m.db.UpdatePluginConf(conf)
}

// RestartPlugin reboots/restart the plugin.
func (m *Manager) RestartPlugin(pluginID string) (string, error) { //TODO update the logic to check if plugin was enabled as it dont work well if it was disabled
	instance, err := m.Instance(pluginID)
	if err != nil {
		return "restart fail", errors.New("instance not found")
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	err = instance.Disable()
	if err != nil {
		return "restart fail", err
	}
	time.Sleep(300 * time.Millisecond)
	err = instance.Enable()
	if err != nil {
		return "restart fail", err
	}
	return "restart ok", nil
}

// PluginInfo returns plugin info.
func (m *Manager) PluginInfo(modulePath string) compat.Info {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if p, ok := m.plugins[modulePath]; ok {
		return p.PluginInfo()
	}
	return compat.Info{
		Name:        "UNKNOWN",
		ModulePath:  modulePath,
		Description: "Oops something went wrong",
	}
}

// Instance returns an instance with the given ID.
func (m *Manager) Instance(pluginID string) (compat.PluginInstance, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if instance, ok := m.instances[pluginID]; ok {
		return instance, nil
	}
	return nil, errors.New("instance not found")
}

// HasInstance returns whether the given plugin ID has a corresponding instance.
func (m *Manager) HasInstance(pluginID string) bool {
	instance, err := m.Instance(pluginID)
	return err == nil && instance != nil
}

type pluginFileLoadError struct {
	Filename        string
	UnderlyingError error
}

func (c pluginFileLoadError) Error() string {
	return fmt.Sprintf("error while loading plugin %s: %s", c.Filename, c.UnderlyingError)
}

func (m *Manager) loadPlugins(directory string) error {
	if directory == "" {
		return nil
	}

	pluginFiles, err := ioutil.ReadDir(directory)
	if err != nil {
		return fmt.Errorf("error while reading directory %s", err)
	}
	for _, f := range pluginFiles {
		pluginPath := filepath.Join(directory, "./", f.Name())
		pRaw, err := plugin.Open(pluginPath)
		if err != nil {
			return pluginFileLoadError{f.Name(), err}
		}
		compatPlugin, err := compat.Wrap(pRaw)
		if err != nil {
			return pluginFileLoadError{f.Name(), err}
		}
		if err := m.LoadPlugin(compatPlugin); err != nil {
			return pluginFileLoadError{f.Name(), err}
		}
	}
	return nil
}

// LoadPlugin loads a compat plugin, exported to sideload plugins for testing purposes.
func (m *Manager) LoadPlugin(compatPlugin compat.Plugin) error {
	modulePath := compatPlugin.PluginInfo().ModulePath
	if _, ok := m.plugins[modulePath]; ok {
		return fmt.Errorf("plugin with module path %s is present at least twice", modulePath)
	}
	m.plugins[modulePath] = compatPlugin
	return nil
}

func (m *Manager) initializePlugins() error {
	for _, p := range m.plugins {
		if err := m.initializePlugin(p); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) initializePlugin(p compat.Plugin) error {
	info := p.PluginInfo()
	instance := p.NewPluginInstance()

	pluginConf, _ := m.db.GetPluginByPath(info.ModulePath)

	if pluginConf == nil {
		var err error
		pluginConf, err = m.createPluginConf(instance, info)
		if err != nil {
			return err
		}
	}

	m.instances[pluginConf.UUID] = instance

	if compat.HasSupport(instance, compat.Storager) {
		instance.SetStorageHandler(dbStorageHandler{pluginConf.UUID, m.db})
	}
	if compat.HasSupport(instance, compat.Configurer) {
		m.initializeConfigurerForSingleUserPlugin(instance, pluginConf)
	}
	if compat.HasSupport(instance, compat.Webhooker) {
		uuid := pluginConf.UUID
		path := pluginConf.ModulePath
		g := m.mux.Group("/"+path, requirePluginEnabled(uuid, m.db))
		instance.RegisterWebhook(strings.Replace(g.BasePath(), ":id", path, 1), g) //change path to uuid if we want the url to register as uuid
	}
	if pluginConf.Enabled {
		err := instance.Enable()
		if err != nil {
			// Single user plugin cannot be enabled
			// Don't panic, disable for now and wait for user to update config
			log.Printf("Plugin initialize failed: %s. Disabling now...", err.Error())
			pluginConf.Enabled = false
			m.db.UpdatePluginConf(pluginConf)
		}
	}
	return nil
}

func (m *Manager) initializeConfigurerForSingleUserPlugin(instance compat.PluginInstance, pluginConf *model.PluginConf) {
	if len(pluginConf.Config) == 0 {
		// The Configurer is newly implemented
		// Use the default config
		pluginConf.Config, _ = yaml.Marshal(instance.DefaultConfig())
		m.db.UpdatePluginConf(pluginConf)
	}
	c := instance.DefaultConfig()
	if yaml.Unmarshal(pluginConf.Config, c) != nil || instance.ValidateAndSetConfig(c) != nil {
		pluginConf.Enabled = false

		log.Printf("Plugin %s failed to initialize because it rejected the current config. It might be outdated. A default config is used and the user would need to enable it again.", pluginConf.ModulePath)
		newConf := bytes.NewBufferString("# Plugin initialization failed because it rejected the current config. It might be outdated.\r\n# A default plugin configuration is used:\r\n")

		d, _ := yaml.Marshal(c)
		newConf.Write(d)
		newConf.WriteString("\r\n")

		newConf.WriteString("# The original configuration: \r\n")
		oldConf := bufio.NewScanner(bytes.NewReader(pluginConf.Config))
		for oldConf.Scan() {
			newConf.WriteString("# ")
			newConf.WriteString(oldConf.Text())
			newConf.WriteString("\r\n")
		}

		pluginConf.Config = newConf.Bytes()

		m.db.UpdatePluginConf(pluginConf)
		instance.ValidateAndSetConfig(instance.DefaultConfig())
	}
}

func (m *Manager) createPluginConf(instance compat.PluginInstance, info compat.Info) (*model.PluginConf, error) {
	pluginConf := &model.PluginConf{
		Name:       info.Name,
		ModulePath: info.ModulePath,
		HasNetwork: info.HasNetwork,
	}
	if compat.HasSupport(instance, compat.Configurer) {
		pluginConf.Config, _ = yaml.Marshal(instance.DefaultConfig())
	}
	if err := m.db.CreatePlugin(pluginConf); err != nil {
		return nil, err
	}
	return pluginConf, nil
}
