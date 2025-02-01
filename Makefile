run:
	@go run main.go --secret=./input/client_secret.json --token=./input/token.json --channel-id=${CHANNEL_ID}

gen:
	@go generate ./...

lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0 run -v

test:
	@go test -v -shuffle=on ./...

tidy:
	@go mod tidy
