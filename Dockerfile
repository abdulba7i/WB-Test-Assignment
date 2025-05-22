FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api/main.go
RUN go build -o consumer ./cmd/consumer/main.go
RUN go build -o publisher ./cmd/publisher/main.go