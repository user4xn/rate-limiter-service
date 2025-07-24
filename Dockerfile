# Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o rate-limiter main.go

EXPOSE 8080

ENTRYPOINT ["./rate-limiter"]
