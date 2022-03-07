package router

import (
	"github.com/GuessWhoSamFoo/octant-synse/pkg/plugin/views"
	"github.com/vmware-tanzu/octant/pkg/plugin/service"
	"github.com/vmware-tanzu/octant/pkg/view/component"
)

func InitRoutes(router *service.Router) {
	router.HandleFunc("", synsePluginListHandler)
	router.HandleFunc("/*", synsePluginDeviceHandler)
}

func synsePluginListHandler(request service.Request) (component.ContentResponse, error) {
	view, err := views.BuildSynsePluginView(request)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	response := component.NewContentResponse(nil)
	response.Add(view)
	return *response, nil
}

func synsePluginDeviceHandler(request service.Request) (component.ContentResponse, error) {
	view, err := views.BuildSynseDeviceView(request)
	if err != nil {
		return component.EmptyContentResponse, err
	}
	response := component.NewContentResponse(nil)
	response.Add(view)
	return *response, nil
}
