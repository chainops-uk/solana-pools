BINARY_NAME=solana-pools

all: build test

build:
	go build -o ./${BINARY_NAME} ./cmd/solana-pools/main.go

test:
	go test

run:
	go build -o ./${BINARY_NAME} ./cmd/solana-pools/main.go
	./${BINARY_NAME} solana pools

build-docker:
	docker build -t ${BINARY_NAME} -f ./Dockerfile .

run-docker:
	docker run -d -p 9861:9861 --name ${BINARY_NAME} --restart unless-stopped --network="host" ${BINARY_NAME} ./${BINARY_NAME} solana pools

stop-container:
	docker stop solana-pools

clean:
	go clean
	rm ${BINARY_NAME}