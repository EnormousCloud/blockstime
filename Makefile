build:
	go build -o ./bin/blockstime-server ./cmd/blockstime-server

run:
	go run ./cmd/blockstime-server

lint:
	golangci-lint run
