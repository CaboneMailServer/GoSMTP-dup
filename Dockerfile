# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o gosmtp-dup .

# Final stage
FROM alpine:latest

# Install ca-certificates and netcat for health check
RUN apk --no-cache add ca-certificates netcat-openbsd

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/gosmtp-dup .

# Copy config file
COPY --from=builder /app/config-example.yaml ./config.yaml

# Create config directory
RUN mkdir -p /etc/smtp-dup

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app /etc/smtp-dup

# Switch to non-root user
USER appuser

# Expose SMTP port
EXPOSE 2525

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD nc -z localhost 2525 || exit 1

# Run the application
CMD ["./gosmtp-dup"]