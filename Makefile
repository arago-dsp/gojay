.PHONY: all dep build test fuzz upgrade help

all: build test

dep: ## Get go.mod dependencies
	@go mod download
	@go mod verify

build: dep ## Build the binary file
	@echo "Building gojay"
	@GOOS=linux GOARCH=amd64 go build ./...

test: ## Run test with a race detector and code coverage report
	# Run the tests and code coverage analysis
	@go test -v -race -covermode=atomic -coverpkg=./...

fuzz: ## Run fuzzing test
	@go test -fuzz=FuzzUnmarshalRaw -fuzztime 30s
	@go test -fuzz=FuzzUnmarshalFields -fuzztime 30s

upgrade: ## Upgrade all the dependencies
	@go get -u -t ./...
	@go mod tidy

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
