# Gebruik de officiële Go-image
FROM golang:1.20

WORKDIR /app

COPY config.json .

COPY go.mod ./
RUN go mod tidy

