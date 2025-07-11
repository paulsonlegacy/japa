# Start from Go base image
FROM golang:1.23.11-alpine

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first
# Then install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project files
COPY . .

# Build the Go app;
RUN go build -o main ./cmd

# Expose the app port
EXPOSE 8080

# Start the app
CMD ["./main"]