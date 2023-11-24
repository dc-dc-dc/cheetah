.PHONY: test coverage example_exchange example_market

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test:
	go test -v ./...

example_exchange:
	go run ./cmd/example_exchange/main.go

example_market:
	go run ./cmd/example_market/main.go