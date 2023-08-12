help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## golangci linting
	golangci-lint run ./...

.PHONY: test
test: ## runs tests
	go test -v -count=1 ./...

test-race: ## test with race detector
	go test -race -count=1 ./...

test-coverage: ## check coverage in default browser
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

run: ## run locally in development mode
	go run main.go
