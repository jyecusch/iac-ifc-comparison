package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, event events.SNSEvent) error {
	for _, record := range event.Records {
		snsRecord := record.SNS
		log.Printf("Received SNS message [ID: %s]: %s", snsRecord.MessageID, snsRecord.Message)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}
