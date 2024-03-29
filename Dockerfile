# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.12 base image
FROM golang:1.12

# Force the go compiler to use modules
ENV GO111MODULE=on

# Add Maintainer Info
LABEL maintainer="Cristian Chaparro <cristianchaparroa@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/cristianchaparroa/secret

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
RUN go mod download
RUN go mod tidy

RUN go build .


# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./secret"]
