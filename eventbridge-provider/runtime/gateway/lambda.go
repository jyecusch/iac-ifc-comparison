package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/jyecusch/iac-ifc-comparison/eventbridge-provider/eventbridge"
	lambda_gateway "github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/topics"
	"google.golang.org/protobuf/types/known/structpb"
)

func handleEventBridgeEvent(_ context.Context, subscriptions topics.SubscriptionRequestHandler, evt json.RawMessage) (interface{}, error) {
	cwEvt := &events.CloudWatchEvent{}
	if err := json.Unmarshal(evt, cwEvt); err != nil {
		return nil, lambda_gateway.NewUnhandledLambdaEventError(err)
	}

	topicName, err := eventbridge.GetTopicNameFromSourceId(cwEvt.Source)
	if err != nil {
		return nil, err
	}

	var message map[string]interface{}

	if err := json.Unmarshal(cwEvt.Detail, &message); err != nil {
		return nil, fmt.Errorf("unable to unmarshal nitric message from EventBridge trigger: %w", err)
	}

	structMessage, err := structpb.NewStruct(message)
	if err != nil {
		return nil, fmt.Errorf("unable to convert message to struct: %w", err)
	}

	topicMessage := &topicspb.TopicMessage{
		Content: &topicspb.TopicMessage_StructPayload{
			StructPayload: structMessage,
		},
	}

	request := &topicspb.ServerMessage{
		Content: &topicspb.ServerMessage_MessageRequest{
			MessageRequest: &topicspb.MessageRequest{
				TopicName: topicName,
				Message:   topicMessage,
			},
		},
	}

	resp, err := subscriptions.HandleRequest(request)
	if err != nil {
		return nil, err
	}

	if !resp.GetMessageResponse().Success {
		return nil, fmt.Errorf("event processing failed")
	}

	return nil, nil
}

func RouteEvent(ctx context.Context, resolver resource.AwsResourceResolver, handlers *lambda_gateway.Handlers, evt json.RawMessage) (interface{}, error) {
	response, err := lambda_gateway.StandardEventRouter(ctx, resolver, handlers, evt)

	if err != nil {
		var unhandledErr *lambda_gateway.UnhandledLambdaEventError

		// Extend the LambdaGateway to handle EventBridge events
		if errors.As(err, &unhandledErr) {
			return handleEventBridgeEvent(ctx, handlers.Subscriptions, evt)
		}
	}

	return response, err
}
