# Build stage
FROM golang:1.24-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Set environment variables for static build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# Focusing on the worker entry point
RUN go build -o worker-app worker/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls and tzdata for timezones
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/worker-app .

# Ensure /tmp exists (it usually does, but explicitly setting permissions if needed is good practice)
# Alpine's /tmp is usually world-writable, which is what we need.

# Command to run the executable
CMD ["./worker-app"]
