FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/events ./cmd/events

FROM alpine:3.20

RUN adduser -D -g '' appuser
COPY --from=builder /out/events /usr/local/bin/events

USER appuser
ENTRYPOINT ["/usr/local/bin/events"]