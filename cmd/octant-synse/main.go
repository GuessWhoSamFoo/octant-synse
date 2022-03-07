package main

import (
	"github.com/GuessWhoSamFoo/octant-synse/pkg/plugin/settings"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
)

func main() {
	name := settings.GetName()
	description := settings.GetDescription()
	capabilities := settings.GetCapabilities()
	options := settings.GetOptions()
	plugin, err := service.Register(name, description, capabilities, options...)
	if err != nil {
		panic(err)
	}
	plugin.Serve()
}
