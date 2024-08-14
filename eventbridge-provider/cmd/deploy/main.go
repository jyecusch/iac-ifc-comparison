package main

import (
	_ "embed"

	"github.com/jyecusch/iac-ifc-comparison/eventbridge-provider/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
)

// Embed the runtime provider
//
//go:embed runtime-extension
var runtimeBin []byte

var runtimeProvider = func() []byte {
	return runtimeBin
}

// Create and start the deployment server
func main() {
	stack := deploy.NewAwsExtensionProvider()

	providerServer := provider.NewPulumiProviderServer(stack, runtimeProvider)

	providerServer.Start()
}
