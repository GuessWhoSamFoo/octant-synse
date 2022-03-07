package synse

import (
	"context"
	"fmt"
	"github.com/vapor-ware/synse-cli/pkg/config"
	"github.com/vapor-ware/synse-cli/pkg/utils"
	"github.com/vapor-ware/synse-client-go/synse"
	"github.com/vapor-ware/synse-client-go/synse/scheme"
	server "github.com/vapor-ware/synse-server-grpc/go"
	"io"
	v1 "k8s.io/api/core/v1"
)

func GetStatus(address string) (*scheme.Status, error) {
	synseClient, err := synse.NewHTTPClientV3(&synse.Options{
		Address: address,
		TLS:     synse.TLSOptions{},
	})
	if err != nil {
		return nil, err
	}
	resp, err := synseClient.Status()
	if err != nil {
		return nil, err
	}
	return resp, err
}

func GetDevices(ctx context.Context, address string) ([]*server.V3Device, error) {
	defer config.Purge()
	devices := new([]*server.V3Device)
	if err := config.AddContext(&config.ContextRecord{
		Name: "tempCtx",
		Type: "plugin",
		Context: config.Context{
			Address: address,
		},
	}); err != nil {
		return *devices, err
	}

	conn, client, err := utils.NewSynseGrpcClient("tempCtx", "")
	if err != nil {
		return *devices, err
	}
	defer conn.Close()

	stream, err := client.Devices(ctx, &server.V3DeviceSelector{})
	if err != nil {
		return *devices, err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return *devices, err
		}
		*devices = append(*devices, resp)
	}

	return *devices, nil
}

func GetAddress(service v1.Service) string {
	var address string
	for _, ing := range service.Status.LoadBalancer.Ingress {
		for _, p := range service.Spec.Ports {
			address = ing.IP + ":" + fmt.Sprintf("%v", p.Port)
			break
		}
	}
	return address
}

func GetReadings(ctx context.Context, device *server.V3Device, address string) ([]*server.V3Reading, error) {
	var readings []*server.V3Reading
	defer config.Purge()
	if err := config.AddContext(&config.ContextRecord{
		Name: "tempCtx",
		Type: "plugin",
		Context: config.Context{
			Address: address,
		},
	}); err != nil {
		return readings, err
	}

	conn, client, err := utils.NewSynseGrpcClient("tempCtx", "")
	if err != nil {
		return readings, err
	}
	defer conn.Close()
	stream, err := client.Read(ctx, &server.V3ReadRequest{Selector: &server.V3DeviceSelector{Id: device.Id}})
	if err != nil {
		return readings, nil
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return readings, nil
		}
		readings = append(readings, resp)
	}
	return readings, nil
}
