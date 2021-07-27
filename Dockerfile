# Set the base image as the official Go image that already has all the tools and packages to compile and run a Go application
FROM registry.ci.openshift.org/open-cluster-management/builder:go1.16-linux AS builder

# Set the current working directory inside the container
WORKDIR /app

COPY . .

# Fetch dependencies.
# Using go get requires root.
USER root
RUN go get -d -v

# Compile the application
RUN go build -o /hello-world-service

# Execute the aplication
CMD [ "/hello-world-service" ]
