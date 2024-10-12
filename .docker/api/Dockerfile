# STEP 1: Build executable sever and purge binaries with UPX (it is an advanced executable file compressor)
FROM golang:1.23 AS builder
# Add Maintainer Info
LABEL maintainer="Bernardo Secades <bernardosecades@gmail.com>"

RUN useradd appuser

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/bernardosecades/zaslink
# Copy go mod and sum files
COPY ../../go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY ../.. .
# Build the Go app

# -ldflags="-w -s" reduce size of bnary
RUN cd cmd/api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/api-secret .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/bin/api-secret"]

# STEP 2: Build a small image
FROM alpine:3.13 AS production

# passwd file to appuser created in first stage
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /bin/api-secret /go/bin/api-secret

RUN  ls -la /go/bin/

# Use an unprivileged user.
USER appuser
# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/go/bin/api-secret"]