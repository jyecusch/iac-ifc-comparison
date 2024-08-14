package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSnsClient is a mock implementation of the SnsClient interface
type MockSnsClient struct {
	mock.Mock
}

func (m *MockSnsClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, params, optFns)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestHandleRequestSuccess(t *testing.T) {
	// Set up environment variables
	os.Setenv("SNS_TOPIC_ARN", "arn:aws:sns:us-west-2:123456789012:MyTopic")

	mockSnsClient := new(MockSnsClient)
	mockSnsClient.On("Publish", mock.Anything, mock.AnythingOfType("*sns.PublishInput"), mock.Anything).Return(&sns.PublishOutput{}, nil)

	handler := &ApiHandler{
		snsClient: mockSnsClient,
	}

	request := events.APIGatewayProxyRequest{
		Body: "Test Message",
	}

	response, err := handler.handleRequest(context.Background(), request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "Message published", response.Body)
	assert.Equal(t, "text/plain; charset=utf-8", response.Headers["Content-Type"])

	mockSnsClient.AssertExpectations(t)
}

func TestHandleRequestError(t *testing.T) {
	// Set up environment variables
	os.Setenv("SNS_TOPIC_ARN", "arn:aws:sns:us-west-2:123456789012:MyTopic")

	mockSnsClient := new(MockSnsClient)
	mockSnsClient.On("Publish", mock.Anything, mock.AnythingOfType("*sns.PublishInput"), mock.Anything).Return(nil, errors.New("failed to publish message"))

	handler := &ApiHandler{
		snsClient: mockSnsClient,
	}

	request := events.APIGatewayProxyRequest{
		Body: "Test Message",
	}

	response, err := handler.handleRequest(context.Background(), request)

	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)
	assert.Equal(t, "Failed to publish message", response.Body)

	mockSnsClient.AssertExpectations(t)
}

func TestHandleRequestMissingTopicArn(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("SNS_TOPIC_ARN")

	mockSnsClient := new(MockSnsClient)

	handler := &ApiHandler{
		snsClient: mockSnsClient,
	}

	request := events.APIGatewayProxyRequest{
		Body: "Test Message",
	}

	_, err := handler.handleRequest(context.Background(), request)
	assert.Error(t, err)
}

func TestNewHandlerSuccess(t *testing.T) {
	os.Setenv("AWS_REGION", "us-west-2")
	os.Setenv("SNS_TOPIC_ARN", "arn:aws:sns:us-west-2:123456789012:MyTopic")

	handler, err := NewHandler()
	assert.NoError(t, err)
	assert.NotNil(t, handler)
	assert.NotNil(t, handler.snsClient)
}

func TestNewHandlerMissingRegion(t *testing.T) {
	os.Unsetenv("AWS_REGION")

	handler, err := NewHandler()
	assert.Error(t, err)
	assert.Nil(t, handler)
}
