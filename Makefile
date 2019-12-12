.PHONY: build clean docker vendor

default: build-local

all: vendor build docker

build: build-darwin build-linux

build-local:
	gofmt -w ./server.go
	gofmt -w ./pkg/domain/*.go
	gofmt -w ./pkg/handler/*.go
	go build -o cookiemonster2 -v ./server.go

build-darwin:
	GOOS=darwin CGO_ENABLED=0 go build -o cookiemonster-darwin-amd64 -v ./server.go

build-linux:
	GOOS=linux CGO_ENABLED=0 go build -o cookiemonster-linux-amd64 -v ./server.go

clean:
	rm -rf ./bin ./vendor

docker:
	docker build --no-cache -t seungkyua/cookiemonster -f Dockerfile.cookiemonster .

