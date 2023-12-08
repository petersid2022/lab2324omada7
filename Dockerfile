# Use an official Golang runtime as a parent image
FROM golang:1.21.4-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any dependencies
RUN go mod download

# Install Air for hot-reloading
RUN go install github.com/cosmtrek/air@latest

# For debugging
RUN ls -al

# Expose port 8080
EXPOSE 8080

# Run the application using Air for hot-reloading
CMD ["/go/bin/air"]
