# my-api-lambda

This is a proof-of-concept API lambda function which expected to be called from an API gateway endpoint.

The data flow is:
* RESTful json data > AWS API Gateway > AWS API Gateway mapping > customEvent > lambda handler

## Deploy Environment Variables

* ENV
    + deployment environment of: prod, qa
* LOG_LEVEL
    + log level of: debug, info, warn, error

## How to deploy to AWS (GUI AWS console)

* build the lambda
    + invoke "make release"
* upload the lambda zip file
    + go to “Lambda > Functions > my-api-lambda > Code”
    + click "Upload from > .zip file", then click "Upload" again and select the built zip file my-api-lambda.zip
    + upload the same file to other lambda function my-api-lambda-healthcheck and my-api-lambda-qa

## How to test in local development

* set environment variables as above accordingly
* invoke 'go run *.go'

## How to test in AWS console

* go to AWS Lambda > Functions
* click on "my-api-lambda" or "my-api-lambda-healthcheck"
* under Test tab, input the following to "Event JSON"
    ```JSON
    {
        "method": "PUT",
        "path": "/health",
        "body": {
            "callback": "ValueChanged",
            "data": {
                "value": 1001
            }
        }
    }
    ```
* click "Test" button should have "Executing function: succeeded" with response
    ```JSON
    {
        "status": "ok",
        "method": "PUT",
        "path": "/callback/company-1001",
        "body": {
            "callback": "ValueChanged",
            "data": {
                "value": 1001
            }
        }
    }
    ```

## API Gateway

The above lambda is supposed to be called from API Gateway.

### How to deploy
* in "API Gateway" > "APIs", click "Create API" button
* in "API Gateway" > "APIs" > "Resources", click "Create resource" and "Create method" accordingly to create a new resource of "PUT" "/callback/company-1001"
* once all APIs are created, click "Deploy API" button
    + select "No stage" for the very first deployment (somehow, must deploy as "No stage" first before you can create a stage)
    + click "Deploy" button to deploy
* in "API Gateway" > "APIs" > "My Callback API (on2i8y7rxj)" > "Stages"
    + click "Create stage" button to create a new stage, in "Deployment", select the datetime of the appropriate deploy
    + click "Create stage" button again to really create it    
* added "Resources" > "Integration request" > "Mapping templates"
    ```
    ##  See https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html
    ##  This template will pass through all parameters including path, querystring, header, stage variables, and context through to the integration endpoint via the body/payload
    #set($allParams = $input.params())
    {
    "method" : "$context.httpMethod",
    "path" : "$context.resourcePath",
    "body" : $input.body,
    "context" : {
        "account-id" : "$context.identity.accountId",
        "api-id" : "$context.apiId",
        "api-key" : "$context.identity.apiKey",
        "authorizer-principal-id" : "$context.authorizer.principalId",
        "caller" : "$context.identity.caller",
        "cognito-authentication-provider" : "$context.identity.cognitoAuthenticationProvider",
        "cognito-authentication-type" : "$context.identity.cognitoAuthenticationType",
        "cognito-identity-id" : "$context.identity.cognitoIdentityId",
        "cognito-identity-pool-id" : "$context.identity.cognitoIdentityPoolId",
        "http-method" : "$context.httpMethod",
        "stage" : "$context.stage",
        "source-ip" : "$context.identity.sourceIp",
        "user" : "$context.identity.user",
        "user-agent" : "$context.identity.userAgent",
        "user-arn" : "$context.identity.userArn",
        "request-id" : "$context.requestId",
        "resource-id" : "$context.resourceId",
        "resource-path" : "$context.resourcePath"
        }
    }
    ```
* AWS gui console test the API resource with
    ```JSON
    {
        "callback": "ValueChanged",
        "data": {
            "value": 1001
        }
    }
    ```
* terminal test
    + invoke "./script/test-my-api.sh"
    + expect response:
    ```JSON
    "{\"status\":\"ok\",\"method\":\"PUT\",\"path\":\"/callback/company-1001\",\"body\":{\"hello\":\"Hello, my first API lambda!\"}}"
    ```

### How to disable
* need to delete all the API stages
* without stage, API wont be accessible

## References

* [AWS Lambda and Golang](https://blog.stackademic.com/aws-lambda-and-golang-72c191294e82)
* [Migrating AWS Lambda functions from the Go1.x runtime to the custom runtime on Amazon Linux 2](https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/)
