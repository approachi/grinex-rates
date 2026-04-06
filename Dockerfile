FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/server ./cmd/server

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/bin/server ./server
COPY --from=builder /app/config.yaml ./config.yaml

EXPOSE 50051

CMD ["./server"]
