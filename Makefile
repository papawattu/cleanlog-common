.PHONY: all build clean client run run-client watch docker-build

all: build client

deps: 
	@go mod tidy
	@go mod vendor
	@go mod verify
	@go mod download
	@go get ./...
build: deps test
	@go build 

test: 
	@go test -v ./...

coverage:
	@go test -v ./... -coverprofile=tmp/coverage.out
	@go tool cover -html=tmp/coverage.out -o tmp/coverage.html
	@rm tmp/coverage.out
	@open tmp/coverage.html