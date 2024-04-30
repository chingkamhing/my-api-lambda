package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"net/url"
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
	logger.Info("http-request-wrapper", "version", version, "build", build, "user", user)
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
			logger.Error("http-request-wrapper", "method", event.Method, "url", event.Url, "response", response, "error", err, "REQUEST ID", lc.AwsRequestID, "FUNCTION NAME", lambdacontext.FunctionName, "DEADLINE", deadline.String())
		} else {
			logger.Info("http-request-wrapper", "method", event.Method, "url", event.Url, "response", response)
		}
		return response, err
	}
	lambda.Start(myHandler)
}

// cliLambdaHandler is a command line handler for local development purpose
func cliLambdaHandler(eventHandler LambdaEventHandler) {
	flagMethod := flag.String("method", "GET", "http method")
	flagUrl := flag.String("url", "https://api-qa.wahkwong.com.hk/health", "http request url")
	flag.Parse()
	u, err := url.Parse(*flagUrl)
	if err != nil {
		log.Fatalf("Invalid url %q: %v", *flagUrl, err)
	}
	var event = customEvent{
		Method: *flagMethod,
		Url:    u.String(),
	}
	logger.Debug("Event handler", "method", *flagMethod, "url", u.String())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	response, err := eventHandler(ctx, event)
	if err != nil {
		logger.Error("http-request-wrapper", "method", event.Method, "url", event.Url, "response", response, "error", err.Error())
	} else {
		logger.Info("http-request-wrapper", "method", event.Method, "url", event.Url, "response", response)
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
