# References:
# - https://aws.amazon.com/blogs/compute/migrating-aws-lambda-functions-from-the-go1-x-runtime-to-the-custom-runtime-on-amazon-linux-2/

# lambda function name
LAMBDA_NAME ?= my-api-lambda
# deploy environment of: qa, prod
ENV ?= qa

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags="-X 'main.version=$$(git describe --tags --abbrev=0)' -X 'main.build=$$(date '+%Y-%m-%dT%H:%M:%S')' -X 'main.user=$${USER}'" -o dist/bootstrap *.go
	CGO_ENABLED=0 go build -v -ldflags="-X 'main.version=$$(git describe --tags --abbrev=0)' -X 'main.build=$$(date '+%Y-%m-%dT%H:%M:%S')' -X 'main.user=$${USER}'" -o dist/bootstrap.native *.go

.PHONY: release
release: build
	zip --junk-paths dist/$(LAMBDA_NAME).zip dist/bootstrap

.PHONY: deploy
deploy: release
	@if [ "$(ENV)" == "prod" ]; then \
		echo "Deploy $(LAMBDA_NAME) to production..." ; \
		docker run --rm -it -v ~/.aws:/root/.aws -v $${PWD}:/aws amazon/aws-cli lambda update-function-code --function-name $(LAMBDA_NAME) --zip-file fileb:///aws/dist/$(LAMBDA_NAME).zip ; \
	else \
		echo "Deploy $(LAMBDA_NAME) to $(ENV)..." ; \
		docker run --rm -it -v ~/.aws:/root/.aws -v $${PWD}:/aws amazon/aws-cli lambda update-function-code --function-name $(LAMBDA_NAME)-$(ENV) --zip-file fileb:///aws/dist/$(LAMBDA_NAME).zip ; \
	fi

.PHONY: clean
clean:
	-rm dist/bootstrap
	-rm dist/bootstrap.native