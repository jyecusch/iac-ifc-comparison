package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type ApiResponse struct {
	Message string `json:"message"`
	Input   string `json:"input"`
}

type ApiHandler struct {
	snsClient SnsClient
}

// handleRequest is the Lambda function handler
func (a *ApiHandler) handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// The code can't easily communicate that it needs these environment variables set.
	// It relies on the Terraform developer knowing these requirements and manually implementing them.
	topicArn := os.Getenv("SNS_TOPIC_ARN")

	// Build time misconfigurations are tricky to detect - so the code must check the value at runtime.
	if topicArn == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "SNS_TOPIC_ARN environment variable is not set",
			Headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
		}, fmt.Errorf("SNS_TOPIC_ARN environment variable is not set")
	}

	message := map[string]interface{}{
		"data": string(request.Body),
	}

	messageBytes, _ := json.Marshal(message)

	publishInput := &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(string(messageBytes)),
	}

	// The code has no idea whether this operation is permitted or not.
	// If the IAM configuration in the Terraform code is incorrect, the operation will fail.
	// Misconfigured IAM is difficult detect before deploying and running the code.
	_, err := a.snsClient.Publish(ctx, publishInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to publish message",
			Headers: map[string]string{
				"Content-Type": "text/plain; charset=utf-8",
			},
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Message published",
		Headers: map[string]string{
			"Content-Type": "text/plain; charset=utf-8",
		},
	}, nil
}

func NewHandler() (*ApiHandler, error) {
	// AWS Specific configuration and libraries make local integration testing difficult
	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		return nil, fmt.Errorf("AWS_REGION environment variable is not set")
	}

	cfg, sessionError := config.LoadDefaultConfig(context.Background(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	snsClient := sns.NewFromConfig(cfg)

	return &ApiHandler{
		snsClient: snsClient,
	}, nil
}

// Enables mocking of the SNS client for unit tests
type SnsClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func main() {
	handler, err := NewHandler()
	if err != nil {
		log.Fatalf("Error creating new API handler: %v", err)
	}
	lambda.Start(handler.handleRequest)
}
