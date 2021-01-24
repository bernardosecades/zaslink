FROM golang:1.15 as builder
WORKDIR $GOPATH/src/github.com/bernardosecades/sharesecret
COPY ./ .
RUN GOOS=linux GOARCH=amd64 go build -ldflags "-X main.commitHash=$(git rev-parse --short HEAD)"
RUN cp sharesecret /

FROM alpine:latest
COPY --from=builder /sharesecret /
CMD ["/sharesecret"]