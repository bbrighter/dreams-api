ARG GO_VERSION="1.23"
ARG ALPINE_VERSION="3.21"
ARG APP_NAME="dreams-api"


FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG TARGETOS
ARG TARGETARCH
ARG APP_NAME

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o ${APP_NAME} .


FROM --platform=$BUILDPLATFORM alpine:${ALPINE_VERSION}
ARG APP_NAME

COPY --from=builder /app/${APP_NAME} /app/${APP_NAME}
WORKDIR /app
ENTRYPOINT ["/app/dreams-api"]