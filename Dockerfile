FROM alpine:latest

RUN apk add --no-cache git make musl-dev go

ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

LABEL maintainer="Andrey Pugach <apu@everstake.one>"

COPY ./ /solana-pools

WORKDIR /solana-pools

RUN go mod download

RUN go build -o ./solana-pools ./cmd/solana-pools/main.go

EXPOSE 9861
CMD ./solana-pools
#docker build -t solana_pools -f ./Dockerfile .
#docker run -d -p 9861:9861 --name solana_pools --restart unless-stopped --network="brig" solana_pools