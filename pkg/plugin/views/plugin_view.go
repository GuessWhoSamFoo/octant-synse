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
)

func BuildSynsePluginView(request service.Request) (component.Component, error) {
	ctx := request.Context()
	client := request.DashboardClient()

	ul, err := client.List(ctx, store.Key{
		APIVersion: "v1",
		Kind:       "Service",
		Selector: &labels.Set{
			"synse-component": "plugin",
		},
	})

	u, err := client.Get(ctx, store.Key{
		Name:       "synse-server", // TODO: Handle case for multiple servers
		Namespace:  "default",
		APIVersion: "v1",
		Kind:       "Service",
		Selector: &labels.Set{
			"app": "synse-server",
		},
	})
	if err != nil {
		return nil, err
	}
	serverService := &v1.Service{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), serverService); err != nil {
		return nil, err
	}

	var address string
	for _, ing := range serverService.Status.LoadBalancer.Ingress {
		for _, p := range serverService.Spec.Ports {
			address = ing.IP + ":" + fmt.Sprintf("%v", p.Port)
			break
		}
	}

	resp, err := synse.GetStatus(address)
	if err != nil {
		return nil, err
	}

	sections := component.SummarySections{}
	sections.AddText("Name", serverService.Name)
	sections.AddText("Status", resp.Status)
	sections.AddText("Address", address)

	summary := component.NewSummary("Server", sections...)

	table := component.NewTableWithRows(
		"Synse Plugins",
		"There are no plugins",
		component.NewTableCols("Name", "Age"),
		[]component.TableRow{})

	for _, item := range ul.Items {
		tr := component.TableRow{
			"Name": component.NewLink(item.GetName(), item.GetName(), item.GetName()),
			"Age":  component.NewTimestamp(item.GetCreationTimestamp().Time),
		}
		table.Add(tr)
	}
	table.Sort("Name")

	fl := component.NewFlexLayout("")
	fl.AddSections(component.FlexLayoutSection{
		{
			Width: component.WidthFull,
			View:  summary,
		},
		{
			Width: component.WidthFull,
			View:  table,
		},
	})
	return fl, nil
}
