# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o lazyservice ./cmd/lazyservice

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata docker-cli kubectl

# Create non-root user
RUN addgroup -g 1001 -S lazyservice && \
    adduser -u 1001 -S lazyservice -G lazyservice

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/lazyservice .

# Copy config files if they exist
COPY --from=builder /app/config* ./config/ 2>/dev/null || true

# Change ownership
RUN chown -R lazyservice:lazyservice /app

# Switch to non-root user
USER lazyservice

# Expose any ports if needed (none for this TUI app)
# EXPOSE 8080

# Health check (optional)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep lazyservice || exit 1

# Run the application
ENTRYPOINT ["./lazyservice"]
