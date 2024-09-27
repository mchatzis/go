FROM golang:1.23.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG BUILD_VERSION
RUN go build -o main -ldflags "-X main.Version=${BUILD_VERSION}" ./cmd/producer
