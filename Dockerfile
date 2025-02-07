# Use the official Golang image as the base image for building
FROM golang:1.20-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download the dependencies to cache them (will be reused if go.mod/go.sum didn't change)
RUN go mod tidy

# Copy the entire Go project to the container
COPY . .

# Build the Go application (this will output the binary as "main")
RUN go build -o main ./main.go

# Start a new stage from a smaller image for running the application
FROM alpine:latest

# Install required certificates for HTTPS communication (optional but often needed)
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go binary from the builder image
COPY --from=builder /app/main .

# Expose the port the app will be running on (e.g., 8080)
EXPOSE 8080

# Command to run the Go binary
CMD ["./main"]
