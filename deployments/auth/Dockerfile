FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/auth ./cmd/auth

FROM alpine:3.20

RUN adduser -D -g '' appuser
COPY --from=builder /out/auth /usr/local/bin/auth
RUN chmod +x /usr/local/bin/auth

USER appuser
ENTRYPOINT ["/usr/local/bin/auth"]