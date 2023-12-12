# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Prefetch the binaries, so that they will be cached and not downloaded on each change
RUN go run github.com/steebchen/prisma-client-go prefetch

# Copy the rest of the application code to the container
COPY . .

# Push migrations to production database
# RUN go run github.com/steebchen/prisma-client-go migrate deploy

# Build the Go application
RUN go build -o main .

# Expose the port that the application will run on
EXPOSE 9000

# Command to run the application
CMD ["./main"]
