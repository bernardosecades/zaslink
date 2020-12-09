FROM golang:latest

# Set the Current Working Directory inside the container
#WORKDIR $GOPATH/src/github.com/bernardosecades/sharesecret

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
#COPY . .

# Download all the dependencies
#RUN go get -d -v ./...

# Install the package
#RUN go install -v ./...

#RUN apt-get update
#RUN apt-get install vim -y
RUN go get "github.com/go-sql-driver/mysql"