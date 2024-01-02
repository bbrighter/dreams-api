FROM golang:bullseye
# RUN apt-get install gcc-aarch64-linux-gnu libc6-dev-arm64-cross PROBABLY NOT NEEDED????
RUN apt update
RUN apt install -y gcc make gcc-arm-linux-gnueabi binutils-arm-linux-gnueabi

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
COPY *.go ./

WORKDIR /dream-api
# env GOARCH=arm GOARM=6 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc  go build .
