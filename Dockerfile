# GO 1.20
FROM golang:1.20-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the container
COPY go.mod go.sum ./

# Download and install any dependencies
RUN go mod download

# Copy the rest of the application code to the container
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port that the application will run on
EXPOSE 3000

# Define environment variable(s)
ARG DB_URL
ARG SUPABASE_URL
ARG SUPABASE_KEY

ENV DB_URL=${DB_URL}
ENV SUPABASE_URL=${SUPABASE_URL}
ENV SUPABASE_KEY=${SUPABASE_KEY}

# Command to run the application
CMD ["./main"]
