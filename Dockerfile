FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o api ./cmd/api

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./api"]
