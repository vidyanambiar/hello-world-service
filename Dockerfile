# Set the base image as the official Go image that already has all the tools and packages to compile and run a Go application
FROM golang:1.16-alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy modules into the working directory
COPY go.mod .
# Install modules within the working directory in the container
RUN go mod download

# Copy source code into the image
COPY *.go .

# Compile the application
RUN go build -o /hello-world-service

# Execute the aplication
CMD [ "/hello-world-service" ]
