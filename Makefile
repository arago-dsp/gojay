.PHONY: all get dep build test fuzz upgrade help

all: build test

get: ## Get latest dependencies
	@go get ./...

dep: ## Get go.mod dependencies
	@go mod download

build: dep ## Build the binary file
	@echo "Building gojay"
	@GOOS=linux GOARCH=amd64 go build ./...

test: ## Run test with a race detector and code coverage report
	# Run the tests and code coverage analysis
	@go test -json -race -run=^Test -covermode=atomic -coverpkg=./... -coverprofile .ci/coverage.txt > .ci/tests.jsonl; RET_CODE=$$?; gotestsum --junitfile .ci/tests.xml --raw-command cat .ci/tests.jsonl; exit $${RET_CODE}
	@go tool cover -func .ci/coverage.txt

fuzz: ## Run fuzzing test
	@go test -fuzz=FuzzUnmarshalRaw -fuzztime 30s
	@go test -fuzz=FuzzUnmarshalFields -fuzztime 30s

upgrade: get build test ## Upgrade all the dependencies
	@go get -u -t ./...
	@go mod tidy

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
