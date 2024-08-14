package deploy

import (
	"github.com/nitrictech/nitric/cloud/aws/deploy"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
)

type EventBridgeTopic struct {
	bus  *cloudwatch.EventBus
	rule *cloudwatch.EventRule
}

type AwsExtendedProvider struct {
	deploy.NitricAwsPulumiProvider

	Topics map[string]*EventBridgeTopic
}

func NewAwsExtensionProvider() *AwsExtendedProvider {
	awsProvider := deploy.NewNitricAwsProvider()

	return &AwsExtendedProvider{
		NitricAwsPulumiProvider: *awsProvider,
		Topics:                  make(map[string]*EventBridgeTopic),
	}
}
