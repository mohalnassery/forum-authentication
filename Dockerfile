# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /server

# Copy the Go module files and download dependencies
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy the rest of the server code
COPY server /server

# Build the Go application
RUN go build -o main .

# Copy the client code (if needed)
COPY client /client

# Verify that the main file exists
RUN ls -l /server

# Set the entrypoint to the built executable
CMD ["./main"]
