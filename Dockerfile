# Step 1: Build stage
FROM golang:1.21-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go Modules manifests (go.mod and go.sum)
COPY go.mod go.sum ./

# Download and install dependencies
RUN go mod tidy

# Copy the rest of the application source code into the container
COPY . .

# Build the Go application (output binary is "shop-backend")
RUN go build -o shop-backend .






# Step 2: Create a minimal runtime image
FROM alpine:latest

# Set the working directory in the container
WORKDIR /root/

# Copy the binary from the build stage
COPY --from=builder /app/shop-backend .

# Copy the configuration file
COPY config/ /root/config/ 
COPY db/ /root/db/


# # Copy the favicon.ico file (if applicable)
# COPY static/ /root/static/

# Install dependencies (e.g., PostgreSQL client)
RUN apk --no-cache add postgresql-client

# Expose the port the application will run on
EXPOSE 8080

# Run the application
CMD ["./shop-backend"]
