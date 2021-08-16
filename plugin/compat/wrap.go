package compat

import (
	"errors"
	"fmt"
	"plugin"

	papiv1 "github.com/NubeDev/flow-framework/plugin/plugin-api"
)

// Wrap wraps around a raw go plugin to provide typesafe access.
func Wrap(p *plugin.Plugin) (Plugin, error) {
	getInfoHandle, err := p.Lookup("GetFlowPluginInfo")
	if err != nil {
		return nil, errors.New("missing GetFlowPluginInfo symbol")
	}
	switch getInfoHandle := getInfoHandle.(type) {
	case func() papiv1.Info:
		v1 := PluginV1{}

		v1.Info = getInfoHandle()
		newInstanceHandle, err := p.Lookup("NewFlowPluginInstance")
		if err != nil {
			return nil, errors.New("missing NewFlowPluginInstance symbol")
		}
		constructor, ok := newInstanceHandle.(func(ctx papiv1.UserContext) papiv1.Plugin)
		if !ok {
			return nil, fmt.Errorf("NewFlowPluginInstance signature mismatch, func(ctx plugin.UserContext) plugin.Plugin expected, got %T", newInstanceHandle)
		}
		v1.Constructor = constructor
		return v1, nil
	default:
		return nil, fmt.Errorf("unknown plugin version (unrecogninzed GetFlowPluginInfo signature %T)", getInfoHandle)
	}
}
