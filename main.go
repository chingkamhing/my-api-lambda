package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

type LambdaHandler func(handler LambdaEventHandler)
type LambdaEventHandler func(ctx context.Context, event customEvent) (string, error)

// the link-time variables that will be overwritten after go build
var version string = "0.0.0"
var build string = "development"
var user string = "nobody"

var logger *slog.Logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getEnvLogLevel(os.Getenv("LOG_LEVEL"))}))
var env = os.Getenv("ENV")

func main() {
	logger.Info("my-api-lambda", "version", version, "build", build, "user", user)
	var lambdaHandler LambdaHandler
	if env == "prod" || env == "qa" {
		// aws environment
		lambdaHandler = LambdaHandler(awsLambdaHandler)
	} else {
		// cli environment
		lambdaHandler = LambdaHandler(cliLambdaHandler)
	}
	lambdaHandler(LambdaEventHandler(HandleLambdaEvent))
}

// awsLambdaHandler is a lambda handler to handle aws event handler
func awsLambdaHandler(eventHandler LambdaEventHandler) {
	myHandler := func(ctx context.Context, event customEvent) (string, error) {
		logger.Debug("AWS environment", "EVENT", event, "REGION", os.Getenv("AWS_REGION"))
		response, err := eventHandler(ctx, event)
		if err != nil {
			lc, _ := lambdacontext.FromContext(ctx)
			deadline, _ := ctx.Deadline()
			logger.Error("my-api-lambda", "method", event.Method, "path", event.Path, "body", event.Body, "response", response, "error", err, "REQUEST ID", lc.AwsRequestID, "FUNCTION NAME", lambdacontext.FunctionName, "DEADLINE", deadline.String())
		} else {
			logger.Info("my-api-lambda", "method", event.Method, "path", event.Path, "body", event.Body, "response", response)
		}
		return response, err
	}
	lambda.Start(myHandler)
}

// cliLambdaHandler is a command line handler for local development purpose
func cliLambdaHandler(eventHandler LambdaEventHandler) {
	flagMethod := flag.String("method", "GET", "http method")
	flagPath := flag.String("path", "/health", "http request path")
	flagBody := flag.String("body", `{"callback": "ValueChanged", }`, "http request path")
	flag.Parse()
	var event = customEvent{
		Method: *flagMethod,
		Path:   *flagPath,
		Body:   *flagBody,
	}
	logger.Debug("Event handler", "method", *flagMethod, "path", *flagPath, "body", *flagBody)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, err := eventHandler(ctx, event)
	if err != nil {
		logger.Error("my-api-lambda", "method", event.Method, "path", event.Path, "body", event.Body, "response", response, "error", err.Error())
	} else {
		logger.Info("my-api-lambda", "method", event.Method, "path", event.Path, "body", event.Body, "response", response)
	}
}

// getEnvLogLevel get the slog level base on env variable of: debug, info, warn, error
func getEnvLogLevel(env string) slog.Level {
	var logLevel slog.Level
	switch env {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	case "info":
		fallthrough
	default:
		logLevel = slog.LevelInfo
	}
	return logLevel
}
