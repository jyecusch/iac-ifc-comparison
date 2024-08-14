package main

import (
	"github.com/jyecusch/iac-ifc-comparison/eventbridge-provider/runtime/gateway"
	"github.com/jyecusch/iac-ifc-comparison/eventbridge-provider/runtime/topic"
	"github.com/nitrictech/nitric/cloud/aws/runtime"
	lambda_gateway "github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	"github.com/nitrictech/nitric/core/pkg/logger"
	"github.com/nitrictech/nitric/core/pkg/server"
)

func main() {

	resolver, err := resource.New()
	if err != nil {
		logger.Fatalf("could not create aws resource resolver: %v", err)
		return
	}

	extGatewayPlugin := lambda_gateway.New(resolver, lambda_gateway.WithRouter(gateway.RouteEvent))
	extTopicPlugin, err := topic.New()
	if err != nil {
		logger.Fatalf("could not create topic plugin: %v", err)
	}

	m, err := runtime.NewAwsRuntimeServer(resolver, server.WithGatewayPlugin(extGatewayPlugin), server.WithTopicsPlugin(extTopicPlugin))
	if err != nil {
		logger.Fatalf("there was an error initializing the AWS runtime server: %v", err)
	}

	server.Run(m)
}
