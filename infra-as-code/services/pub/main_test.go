package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventBridgeClient is a mock implementation of the EventBridgeClient interface
type MockEventBridgeClient struct {
	mock.Mock
}

func (m *MockEventBridgeClient) PutEvents(ctx context.Context, params *eventbridge.PutEventsInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*eventbridge.PutEventsOutput), args.Error(1)
}

func TestHandleRequestSuccess(t *testing.T) {
	// Set up environment variables
	os.Setenv("EVENT_BUS_NAME", "testbus")
	os.Setenv("EVENT_SOURCE_ID", "testsource")

	mockEventBridgeClient := new(MockEventBridgeClient)
	mockEventBridgeClient.On("PutEvents", mock.Anything, mock.AnythingOfType("*eventbridge.PutEventsInput"), mock.Anything).Return(&eventbridge.PutEventsOutput{}, nil)

	handler := &ApiHandler{
		eventbridgeClient: mockEventBridgeClient,
	}

	request := events.APIGatewayProxyRequest{
		Body: "Test Message",
	}

	response, err := handler.handleRequest(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Message published", response.Body)
	assert.Equal(t, "text/plain; charset=utf-8", response.Headers["Content-Type"])

	mockEventBridgeClient.AssertExpectations(t)
}

func TestHandleRequestError(t *testing.T) {
	// Set up environment variables
	os.Setenv("EVENT_BUS_NAME", "testbus")
	os.Setenv("EVENT_SOURCE_ID", "testsource")

	mockEventBridgeClient := new(MockEventBridgeClient)
	mockEventBridgeClient.On("PutEvents", mock.Anything, mock.AnythingOfType("*eventbridge.PutEventsInput"), mock.Anything).Return(nil, errors.New("failed to publish message"))

	handler := &ApiHandler{
		eventbridgeClient: mockEventBridgeClient,
	}

	request := events.APIGatewayProxyRequest{
		Body: "Test Message",
	}

	response, err := handler.handleRequest(context.Background(), request)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "Failed to publish message", response.Body)

	mockEventBridgeClient.AssertExpectations(t)
}

func TestHandleRequestMissingTopicArn(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("EVENT_BUS_NAME")
	os.Unsetenv("EVENT_SOURCE_ID")

	mockEventBridgeClient := new(MockEventBridgeClient)

	handler := &ApiHandler{
		eventbridgeClient: mockEventBridgeClient,
	}

	request := events.APIGatewayProxyRequest{
		Body: "Test Message",
	}

	_, err := handler.handleRequest(context.Background(), request)
	assert.Error(t, err)
}

func TestNewHandlerSuccess(t *testing.T) {
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("EVENT_BUS_NAME", "testbus")
	os.Setenv("EVENT_SOURCE_ID", "testsource")

	handler, err := NewHandler()
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.eventbridgeClient)
}

func TestNewHandlerMissingRegion(t *testing.T) {
	os.Unsetenv("AWS_REGION")

	handler, err := NewHandler()
	assert.Error(t, err)
	assert.Nil(t, handler)
}
