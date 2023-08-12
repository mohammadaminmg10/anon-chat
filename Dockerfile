# Use an official Golang runtime as a parent image
FROM golang:1.20.2

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code and other necessary files
COPY . .

# Download dependencies and build the Go application
RUN go mod download && \
    GOOS=linux GOARCH=amd64 go build -o anon-chat && \
    chmod +x anon-chat

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./anon-chat"]
