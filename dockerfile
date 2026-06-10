# --- Stage 1: Development & Build ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Install air for local development hot-reloading
RUN go install github.com/air-verse/air@latest

# Copy dependency files and download them
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the actual production binary (Static compilation)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# --- Stage 2: Production Deployment ---
FROM alpine:3.20

WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose your application port (change 8080 to your app's port)
EXPOSE 8080

# Run the compiled binary directly (No air, no Go SDK)
CMD ["./main"]