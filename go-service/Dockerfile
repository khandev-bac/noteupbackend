
FROM golang:1.25.7-bookworm AS builder

RUN apt-get update && apt-get install -y \
    git \
    build-essential \
    ca-certificates \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux \
    go build -o app ./cmd/main.go


FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 3000

CMD ["./app"]
