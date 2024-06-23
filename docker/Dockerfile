# Start by using the official Golang image as the build image
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

COPY ./userservice .
COPY ./docker/config.yaml .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Build the Go app
RUN go build -o user-service user-service/cmd/userservice

# Start a new stage from scratch
FROM alpine:latest  

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/user-service .
COPY --from=builder /app/config.yaml .

# Expose port 8080 to the outside world
EXPOSE 80 443

# Command to run the executable
CMD ["./user-service"]
