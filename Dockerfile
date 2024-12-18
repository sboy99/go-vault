# Use the latest stable Go version
FROM golang:1.23-alpine AS build

WORKDIR /app

# Set necessary environment variables
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Copy all files into the container
COPY . .

# Download Go dependencies
RUN go mod download

# Build the Go binary
RUN go build -o go-vault

# Use a lightweight Alpine image for the final container
FROM alpine:3.18

WORKDIR /app

# Add a non-root user
RUN adduser -D user

# Copy the binary from the build stage
COPY --from=build /app/go-vault ./go-vault

# Set ownership of the binary
RUN chown user:user go-vault

# Set the user for running the container
USER user

# Command to run the application
CMD ["./go-vault"]
