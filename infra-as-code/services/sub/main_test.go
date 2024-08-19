package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {

	// Create a mock EventBridge (CloudWatch) event
	event := events.CloudWatchEvent{
		ID:     "1234",
		Detail: []byte(`{"data": "test message"}`),
	}

	// Call the handler function
	err := handleRequest(context.Background(), event)

	// Validate the result
	assert.NoError(t, err, "Expected no error from handleRequest")
}
