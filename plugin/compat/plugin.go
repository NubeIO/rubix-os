package compat

// Plugin is an abstraction of plugin handler.
type Plugin interface {
	PluginInfo() Info
	NewPluginInstance() PluginInstance
	APIVersion() string
}

// Info is the plugin info.
type Info struct {
	Version     string
	Author      string
	Name        string
	Website     string
	Description string
	License     string
	ModulePath  string
	HasNetwork  bool
}

func (c Info) String() string {
	if c.Name != "" {
		return c.Name
	}
	return c.ModulePath
}
