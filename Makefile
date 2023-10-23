run:
	@go run ./cmd/api

test:
	@GOFLAGS="-count=1" go test -v -cover -race ./...
