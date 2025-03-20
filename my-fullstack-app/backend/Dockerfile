# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install swag CLI and dependencies
RUN apk add --no-cache git && \
    go install github.com/swaggo/swag/cmd/swag@latest

# Add the Go bin to PATH to ensure swag is available
ENV PATH="${PATH}:/root/go/bin"

# Copy go mod and sum files
COPY backend/go.mod backend/go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY backend/ .

# Generate swagger docs
RUN swag init -g cmd/server/main.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Final stage
FROM alpine:3.16

# Add necessary packages
RUN apk --no-cache add ca-certificates tzdata postgresql-client

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/server .

# Copy swagger docs
COPY --from=builder /app/docs ./docs

# Copy the migrations
COPY sql/migrations /app/migrations

# Copy entrypoint script
COPY backend/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Expose the application port
EXPOSE 8080

# Use the entrypoint script
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["./server"]