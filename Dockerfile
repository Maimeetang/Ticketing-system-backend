# === Stage 1: Build stage ===
FROM golang:1.26 AS builder

RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/local

# === Stage 2: Final runtime stage ===
FROM golang:1.26-bookworm

WORKDIR /app

COPY --from=builder /go/bin/air /usr/local/bin/air

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
