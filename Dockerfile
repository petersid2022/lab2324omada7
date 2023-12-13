# Use an official Golang Alpine runtime as a parent image
FROM golang:alpine

# Set the working directory to /app
WORKDIR /app

COPY go.mod .

COPY go.sum .

# Only download updates if modules files have changed
RUN go mod download

COPY ./ ./

# Install Air for hot-reloading
RUN go install github.com/cosmtrek/air@latest

# Run the application using Air for hot-reloading
CMD ["/go/bin/air"]
