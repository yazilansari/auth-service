FROM golang:1.25.10-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o server ./cmd/server

RUN go build -o worker ./cmd/worker

# =========================

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/worker .

COPY .env .

EXPOSE 8080