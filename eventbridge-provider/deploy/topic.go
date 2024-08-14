package deploy

import (
	"fmt"

	"github.com/jyecusch/iac-ifc-comparison/eventbridge-provider/eventbridge"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudwatch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *AwsExtendedProvider) Topic(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Topic) error {
	var err error

	a.Topics[name] = &EventBridgeTopic{}

	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	a.Topics[name].bus, err = cloudwatch.NewEventBus(ctx, name, &cloudwatch.EventBusArgs{
		Name: pulumi.String(eventbridge.BusName(a.StackId, name)),
		Tags: pulumi.ToStringMap(tags.Tags(a.StackId, name, resources.Topic)),
	}, opts...)
	if err != nil {
		return err
	}

	// Filter based on the appropriate source, not technically necessary but good practice
	evtPattern := fmt.Sprintf(`{"source": ["%s"]}`, eventbridge.SourceId(a.StackId, name))

	a.Topics[name].rule, err = cloudwatch.NewEventRule(ctx, name, &cloudwatch.EventRuleArgs{
		Name:         pulumi.String(name),
		EventBusName: a.Topics[name].bus.Name,
		EventPattern: pulumi.String(evtPattern),
	}, opts...)
	if err != nil {
		return err
	}

	for _, sub := range config.Subscriptions {
		targetLambda, ok := a.Lambdas[sub.GetService()]
		if !ok {
			return fmt.Errorf("unable to find lambda %s for subscription", sub.GetService())
		}

		err := createSubscription(ctx, parent, name, a.Topics[name], targetLambda)
		if err != nil {
			return err
		}
	}

	return nil
}

func createSubscription(ctx *pulumi.Context, parent pulumi.Resource, name string, topic *EventBridgeTopic, target *lambda.Function) error {
	var err error

	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	_, err = lambda.NewPermission(ctx, name+"Permission", &lambda.PermissionArgs{
		SourceArn: topic.rule.Arn,
		Function:  target.Name,
		Principal: pulumi.String("events.amazonaws.com"),
		Action:    pulumi.String("lambda:InvokeFunction"),
	}, opts...)
	if err != nil {
		return err
	}

	_, err = cloudwatch.NewEventTarget(ctx, name+"Target", &cloudwatch.EventTargetArgs{
		EventBusName: topic.bus.Name,
		Rule:         topic.rule.Name,
		Arn:          target.Arn,
	}, opts...)
	if err != nil {
		return err
	}

	return nil
}
