# Start from Go base image
FROM golang:1.23.11-alpine

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy entire project
COPY . .

# Build the Go app from cmd/main.go, output as 'main'
RUN go build -o main ./cmd/main.go

# Make sure wait-for-it and entrypoint are executable
RUN chmod +x wait-for-it.sh entrypoint.sh

# Entrypoint script runs wait-for-it and then the app
ENTRYPOINT ["./entrypoint.sh"]

# Expose app port (can be dynamic via env)
EXPOSE 8080