FROM golang:latest

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the application's port
EXPOSE 8080

# Run the application
CMD ["./main"]
