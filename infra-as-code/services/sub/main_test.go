package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandleRequest(t *testing.T) {
	// Create a mock SNS event
	event := events.SNSEvent{
		Records: []events.SNSEventRecord{
			{
				SNS: events.SNSEntity{
					Message:   "This is a test SNS message",
					MessageID: "12345",
				},
			},
		},
	}

	// Call the handler function
	err := handleRequest(context.Background(), event)

	// Validate the result
	assert.NoError(t, err, "Expected no error from handleRequest")
}
