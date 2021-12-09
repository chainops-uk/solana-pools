BINARY_NAME=solana-pools

all: build test

build:
	go build -o ./${BINARY_NAME} ./cmd/solana-pools/main.go

test:
	go test

run:
	go build -o ./${BINARY_NAME} ./cmd/solana-pools/main.go
	./${BINARY_NAME} api twitter-report

build-docker:
	docker build -t ${BINARY_NAME} -f ./Dockerfile .

run-docker:
	docker run -d -p 9861:9861 --name solana-pools --restart unless-stopped --network="host" solana-pools

stop-container:
	docker stop solana-pools

clean:
	go clean
	rm ${BINARY_NAME}