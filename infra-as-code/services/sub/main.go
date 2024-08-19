package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event events.CloudWatchEvent) error {
	log.Printf("Received message [ID: %s]: %s", event.ID, string(event.Detail))

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
