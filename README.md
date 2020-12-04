# Structure folder

https://kgolding.co.uk/blog/2020/02/19/golang-application-directory-structure/

./	The root of the git repo
./README.md	The projects main readme
./go.mod	Created by running go mod github.com/kgolding.go-app-structure
./cmd/*	Folders for each build-able application main.go
./internal/*	Folders for each private package (that can not be used in other projects)
./pkg/*	Folders for each public package (that might be used in other projects)
./vendor/*	Optional: External dependencies as populated by go mod vendor