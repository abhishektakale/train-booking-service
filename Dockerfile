# Stage 1: Build the Go binary
FROM golang:1.23 as builder

# Set the working directory inside the builder image
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the binary for Linux
RUN CGO_ENABLED=0 go build -v -o server cmd/server/main.go

# Stage 2: Build the runtime image
FROM debian:bullseye-slim

# Set the working directory in the runtime container
WORKDIR /app

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/server /app/server

# Expose the application port
EXPOSE 7001

# Ensure the binary is executable
RUN chmod +x /app/server

# Command to run the binary
CMD ["/app/server"]