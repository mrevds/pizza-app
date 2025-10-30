# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/app ./cmd/user-service

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /build/app .

# Copy config
COPY config.yaml .

# Expose gRPC port
EXPOSE 50051

# Run the application
CMD ["./app"]

