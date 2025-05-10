# Use the official Go image as a build stage
FROM golang:1.24 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o hello .

# Use a minimal base image for the final container
FROM gcr.io/distroless/base-debian11

# Set the working directory
WORKDIR /

# Copy the compiled binary from the builder stage
COPY --from=builder /app/hello .

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["./hello"]