mkfile_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
.PHONY: test coverage example_exchange example_market


coverage:
	TESTING_DATA_FILE=$(mkfile_path)/testing_data.csv go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test:
	TESTING_DATA_FILE=$(mkfile_path)/testing_data.csv go test -v ./...

example_exchange:
	go run ./cmd/example_exchange/main.go

example_market:
	go run ./cmd/market_example/main.go