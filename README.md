# Structure folder

https://kgolding.co.uk/blog/2020/02/19/golang-application-directory-structure/

./	The root of the git repo
./README.md	The projects main readme
./go.mod	Created by running go mod github.com/kgolding.go-app-structure
./cmd/*	Folders for each build-able application main.go
./internal/*	Folders for each private package (that can not be used in other projects)
./pkg/*	Folders for each public package (that might be used in other projects)
./vendor/*	Optional: External dependencies as populated by go mod vendor

# Gorilla mux

https://medium.com/@hugo.bjarred/rest-api-with-golang-and-mux-e934f581b8b5

API => muy bueno: https://www.soberkoder.com/go-rest-api-gorilla-mux/

# Service encapsulation structure folder

Muy bueno: https://github.com/irahardianto/service-pattern-go

https://irahardianto.github.io/service-pattern-go/

# Docker

Read: https://qiita.com/osk_kamui/items/1539ade3c23f58b89f80

docker-compose up --build
docker exec -it golang_db
docker exec -it golang_app bash -c "go run main.go"

Golang docker and test pipeline.

https://codefresh.io/docs/docs/learn-by-example/golang/golang-hello-world/


# Architecture

https://www.perimeterx.com/tech-blog/2019/ok-lets-go/
