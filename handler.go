package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

type customEvent struct {
	events.CloudWatchEvent
	// custom event fields that is added in AWS Amazon EventBridge schedule > Target > Payload
	Method string `json:"method"`
	Url    string `json:"url"`
}

// HandleLambdaEvent handle custom event which get the url and perform http request accordingly.
func HandleLambdaEvent(ctx context.Context, event customEvent) (string, error) {
	logger.Info("My Callback Lambda", "event", event)
	response := "My dummy callback lambda"
	return response, nil
}
