# Multi-stage build for Sentinel Gateway

# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sentinel-gateway .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/sentinel-gateway .

# Copy configuration files
COPY --from=builder /app/config.yaml ./config.yaml

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./sentinel-gateway"]