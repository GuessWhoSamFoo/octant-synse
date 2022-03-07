package views

import (
	"fmt"
	"github.com/GuessWhoSamFoo/octant-synse/pkg/synse"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/store"
	"github.com/vmware-tanzu/octant/pkg/view/component"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

func BuildSynseDeviceView(request service.Request) (component.Component, error) {
	pluginName := strings.TrimPrefix(request.Path(), "/")

	ctx := request.Context()
	client := request.DashboardClient()

	u, err := client.Get(ctx, store.Key{
		Name:       pluginName,
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Service",
		Selector: &labels.Set{
			"app":             pluginName,
			"synse-component": "plugin",
		},
	})
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, fmt.Errorf("cannot find plugin: %v", pluginName)
	}

	pluginService := &v1.Service{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pluginService); err != nil {
		return nil, err
	}

	address := synse.GetAddress(*pluginService)

	devices, err := synse.GetDevices(ctx, address)
	if err != nil {
		return nil, err
	}

	table := component.NewTableWithRows("Devices",
		"There are no devices",
		component.NewTableCols("ID", "Alias", "Type", "Reading", "Info"),
		[]component.TableRow{})

	for _, device := range devices {
		var current string
		readings, err := synse.GetReadings(ctx, device, address)
		if err != nil {
			return nil, err
		}
		if len(readings) != 0 {
			r := readings[len(readings)-1]

			current = fmt.Sprintf("%d%s", r.GetInt64Value(), r.Unit.Symbol)
		}
		tr := component.TableRow{
			"ID":      component.NewText(device.Id),
			"Alias":   component.NewText(device.Alias),
			"Type":    component.NewText(device.Type),
			"Reading": component.NewText(current),
			"Info":    component.NewText(device.Info),
		}
		table.Add(tr)
	}
	table.Sort("ID")

	fl := component.NewFlexLayout("")
	fl.AddSections(component.FlexLayoutSection{
		{
			Width: component.WidthFull,
			View:  table,
		},
	})
	return fl, nil
}
