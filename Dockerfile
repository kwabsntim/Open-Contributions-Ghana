# Dockerfile for Open Source Ghana Go application
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server ./cmd/server.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Copy web assets
COPY --from=builder /app/web ./web

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server"]
