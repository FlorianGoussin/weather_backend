# Use an official Go runtime as a parent image
FROM golang:1.21.7-bookworm

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code into the container
COPY . .

# Download and install any required dependencies
RUN go mod download

# Build the Go application
RUN go build -o app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app"]