.PHONY: run
run:
	@go run main.go --secret=./input/client_secret.json --token=./input/token.json --channel-id=${CHANNEL_ID}

.PHONY: gen
gen:
	@go generate ./...

.PHONY: lint
lint:
	@go tool github.com/golangci/golangci-lint/cmd/golangci-lint run -v

.PHONY: test
test:
	@go test -v -shuffle=on ./...

.PHONY: tidy
tidy:
	@go mod tidy
