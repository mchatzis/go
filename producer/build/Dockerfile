FROM golang:1.23.1-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/producer

EXPOSE 8081

CMD ./main --loglevel $LOG_LEVEL
