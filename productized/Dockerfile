# Multi-stage build for the DriftLock API

# Build stage
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o driftlock-api ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -s /bin/sh driftlock

# Set working directory
WORKDIR /home/driftlock

# Copy the binary from builder stage
COPY --from=builder /app/driftlock-api .

# Change ownership to the non-root user
RUN chown -R driftlock:driftlock /home/driftlock

# Switch to non-root user
USER driftlock

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./driftlock-api"]