# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.15

# Add Maintainer Info
LABEL maintainer="Bernardo Secades <bernardosecades@gmail.com>"
# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/bernardosecades/sharesecret
# Copy go mod and sum files
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN cd cmd/server && go build -o ./server .
RUN cd cmd/purge && go build -o ./purge .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./cmd/server/server"]