package settings

import (
	"github.com/vmware-tanzu/octant/pkg/plugin"
)

func GetCapabilities() *plugin.Capabilities {
	return &plugin.Capabilities{
		IsModule:    true,
		ActionNames: []string{},
	}
}
