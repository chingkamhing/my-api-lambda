package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

type customEvent struct {
	events.CloudWatchEvent
	// custom event fields that is added in AWS Amazon EventBridge schedule > Target > Payload
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Body   interface{} `json:"body"`
}

type customResponse struct {
	Status string      `json:"status"`
	Method string      `json:"method"`
	Path   string      `json:"path"`
	Body   interface{} `json:"body"`
}

// HandleLambdaEvent handle custom event which get the url and perform http request accordingly.
func HandleLambdaEvent(ctx context.Context, event customEvent) (string, error) {
	logger.Info("My Callback Lambda", "event", event)
	response := customResponse{
		Status: "ok",
		Method: event.Method,
		Path:   event.Path,
		Body:   event.Body,
	}
	responseByte, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(responseByte), nil
}
