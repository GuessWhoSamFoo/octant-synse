## Synse Server/Emulator Demo

This repository contains a proof of concept for a dashboard via synse-sdk.

It shows how synse-cli contexts can be selected implicitly through user interaction and
potential for device management through a Kubernetes control plane.

### Requirments

 - kind
 - Go 1.17+
 - Octant 0.25.1+

1. Create a kind cluster with an installation of [metallb](https://kind.sigs.k8s.io/docs/user/loadbalancer/)
 to allow services to have external IPs.
2. Install synse server and emulator plugin via `kubectl apply -f manifests/synse-demo.yaml`.
3. Compile the plugin with `make build` and move binary to `$HOME/.config/octant/plugins`.
