.PHONY: help dependencies up start stop restart status ps clean and execute tests

up:
	docker-compose up --build
down:
	docker-compose down
test-coverage:
	docker-compose exec service go clean -testcache && go test  ./... -tags=unit,integration,e2e -coverprofile cover.out && go tool cover -html=cover.out
test-all:
	docker-compose exec service go clean -testcache && go test  ./... -tags=unit,integration,e2e
test-unit:
	docker-compose exec service go clean -testcache && go test ./... -tags=unit
test-integration:
	docker-compose exec service go clean -testcache && go test ./... -tags=integration
test-e2e:
	docker-compose exec service go clean -testcache && go test -v ./... -tags=e2e
