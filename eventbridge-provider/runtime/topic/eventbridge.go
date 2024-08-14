package topic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/smithy-go"
	nitriceventbridge "github.com/jyecusch/iac-ifc-comparison/eventbridge-provider/eventbridge"
	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	"github.com/nitrictech/nitric/core/pkg/help"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type EventBridgeTopicService struct {
	eventbridgeClient *eventbridge.Client
}

var _ topicspb.TopicsServer = (*EventBridgeTopicService)(nil)

func (s *EventBridgeTopicService) Publish(ctx context.Context, req *topicspb.TopicPublishRequest) (*topicspb.TopicPublishResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("EventBridgeEventService.Publish")

	jsonBytes, err := protojson.Marshal(req.Message.GetStructPayload())
	if err != nil {
		return nil, newErr(
			codes.Unknown,
			fmt.Sprintf("unable to serialize message. %s", help.BugInNitricHelpText()),
			err,
		)
	}

	if req.Delay != nil && req.Delay.AsDuration() > 0 {
		err = s.publishDelayed(ctx, req.TopicName, req.Delay.AsDuration(), jsonBytes)
	} else {
		err = s.publish(ctx, req.TopicName, jsonBytes)
	}

	if err != nil {
		if isCloudWatchAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to publish to topic, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return nil, newErr(codes.Internal, "error publishing message", err)
	}

	return &topicspb.TopicPublishResponse{}, nil
}

func (s *EventBridgeTopicService) publish(ctx context.Context, topicName string, message json.RawMessage) error {
	stackId := commonenv.NITRIC_STACK_ID.String()

	busName := nitriceventbridge.BusName(stackId, topicName)
	source := nitriceventbridge.SourceId(stackId, topicName)

	// messageJson := fmt.Sprintf(`{"message":"%s"}`, message)

	eventInput := &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				EventBusName: aws.String(busName),
				Source:       aws.String(source),
				Detail:       aws.String(string(message)),
				DetailType:   aws.String("customDetailType"),
			},
		},
	}

	_, err := s.eventbridgeClient.PutEvents(ctx, eventInput)

	return err
}

func (s *EventBridgeTopicService) publishDelayed(_ context.Context, _ string, _ time.Duration, _ json.RawMessage) error {
	return status.Errorf(codes.Unimplemented, "delayed messaging not implemented in this demo provider")
}

func isCloudWatchAccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "CloudWatch" && strings.Contains(opErr.Unwrap().Error(), "AuthorizationError")
	}
	return false
}

// New - Create a new EventBridge Topic Service
func New() (*EventBridgeTopicService, error) {
	awsRegion := os.Getenv("AWS_REGION")

	cfg, sessionError := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	eventbridgeClient := eventbridge.NewFromConfig(cfg)

	return &EventBridgeTopicService{
		eventbridgeClient: eventbridgeClient,
	}, nil
}
