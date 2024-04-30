# wk-cronjob-lambda

This is a proof-of-concept callback lambda function which expected to be called from an API gateway endpoint. 

## Deploy Environment Variables

* ENV
    + deployment environment of: prod, qa
* LOG_LEVEL
    + log level of: debug, info, warn, error

## How to deploy to AWS (GUI AWS console)

* build the lambda
    + invoke "make release"
* upload the lambda zip file
    + go to “Lambda > Functions > my-callback-lambda > Code”
    + click "Upload from > .zip file", then click "Upload" again and select the built zip file my-callback-lambda.zip
    + upload the same file to other lambda function my-callback-lambda-healthcheck and my-callback-lambda-qa

## How to test in local development

* set environment variables as above accordingly
* invoke 'go run *.go'
    + default health check "[api-qa.wahkwong.com.hk](https://api-qa.wahkwong.com.hk/)"
    + no JWT_TOKEN is needed
* invoke 'go run *.go -endpoint shipwatch/v1/health'
    + health check shipwatch data pipeline
    + need to set JWT_TOKEN accordingly

## How to test in AWS console

* go to AWS Lambda > Functions
* click on "my-callback-lambda" or "my-callback-lambda-healthcheck"
* under Test tab, input the following to "Event JSON"
    ```JSON
    {
        "url": "https://api-qa.wahkwong.com.hk/health",
        "method": "GET"
    }
    ```
* click "Test" button should have "Executing function: succeeded" with response
    ```JSON
    {"status":"ok"}
    ```

## References

* [AWS Lambda and Golang](https://blog.stackademic.com/aws-lambda-and-golang-72c191294e82)
* [Migrating AWS Lambda functions from the Go1.x runtime to the custom runtime on Amazon Linux 2](https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/)
