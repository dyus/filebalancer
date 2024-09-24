FROM golang:1.22-alpine as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server ./cmd/server
RUN go build -o storage ./cmd/storage
RUN go build -o client ./cmd/client

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /app/server .
COPY --from=builder /app/storage .
COPY --from=builder /app/client .
