FROM golang:1.22.4-bullseye
# RUN apt-get install gcc-aarch64-linux-gnu libc6-dev-arm64-cross PROBABLY NOT NEEDED????
RUN apt update
RUN apt install -y gcc make gcc-arm-linux-gnueabi binutils-arm-linux-gnueabi

COPY go.mod go.sum ./

RUN go mod download
COPY *.go ./

WORKDIR /dreams-api
# env GOARCH=arm GOARM=6 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc  go build .
