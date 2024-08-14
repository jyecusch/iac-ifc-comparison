test: generate-mocks
	@echo "Running tests"
	@go test -v ./...

generate-mocks:
	@echo "Generating mocks"
	@mkdir -p mocks/handler
	@mkdir -p mocks/nitric
	@go run github.com/golang/mock/mockgen github.com/nitrictech/go-sdk/handler HttpRequest,MessageRequest > mocks/handler/mock.go
	@go run github.com/golang/mock/mockgen github.com/nitrictech/go-sdk/nitric Topic > mocks/nitric/mock.go

run:
	@echo "Running the IfC example for local testing"
	@nitric start

up: build-provider
	@echo "Deploying the IfC example to the cloud"
	@nitric up

down:
	@echo "Tearing down the IfC example"
	@nitric down


# Build the provider with EventBridge support
build-provider:
	@$(MAKE) -C ../eventbridge-provider install

.PHONY: test generate-mocks run up down