# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./cmd/server

# Stage 2: Create the final, smaller image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /main .

# Expose port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]