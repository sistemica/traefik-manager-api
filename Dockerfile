# Stage 1: Build the application
FROM golang:1.23.5-alpine AS builder

# Set working directory
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache git make

# Copy go module files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o traefik-manager ./cmd/server/main.go

# Stage 2: Create the runtime container
FROM alpine:3.17

# Add CA certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create app directories
RUN mkdir -p /app/data /app/config
WORKDIR /app

# Copy only the executable from the builder stage
COPY --from=builder /app/traefik-manager /app/

# Create a directory for persistent data and set permissions
RUN mkdir -p /data && chown -R appuser:appgroup /data /app

# Copy configuration (if needed)
# COPY config/* /app/config/

# Set environment variables
ENV SERVER_PORT=9000 \
    STORAGE_FILE_PATH=/data/traefik-manager.json \
    LOG_LEVEL=info \
    LOG_FORMAT=json

# Expose the API port
EXPOSE 9000

# Switch to non-root user
USER appuser

# Runtime command
ENTRYPOINT ["/app/traefik-manager"]