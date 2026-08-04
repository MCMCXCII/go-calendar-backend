FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o auth ./cmd/auth

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/auth .
COPY config ./config

CMD ["./auth"]