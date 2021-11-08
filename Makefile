BINARY_NAME=solana-pool

all: build test

build:
	go build -o ./${BINARY_NAME} ./cmd/solana-pool/main.go

test:
	go test

run:
	go build -o ./${BINARY_NAME} ./cmd/solana-pool/main.go
	./${BINARY_NAME} api twitter-report

build-docker:
	docker build -t eversol_back -f ./Dockerfile .

run-docker:
	docker run -d -p 8962:8962 --name eversol_back_twitter_report --restart unless-stopped --network="host" eversol_back ./

migration-up:
	go run ./cmd/cli/main.go migration up

clean:
	go clean
	rm ${BINARY_NAME}