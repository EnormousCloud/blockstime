build:
	go build -o ./bin/blockstime-server ./cmd/blockstime-server

run:
	go run ./cmd/blockstime-server

index-mainnet:
	go run ./cmd/blockstime-server -i eth.mainnet 

index-rinkeby: build
	./bin/blockstime-server --index eth.rinkeby

index-goerli:
	go run ./cmd/blockstime-server -i eth.goerli

lint:
	golangci-lint run
