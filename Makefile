.PHONY: build clean docker

default: build-local

all: build docker

build: build-darwin build-linux

build-local:
	gofmt -w ./src/cmd/*.go
	gofmt -w ./src/domain/*.go
	gofmt -w ./src/handler/*.go
	go build -o bin/cookiemonster2 -v ./src/cmd

build-darwin:
	GOOS=darwin CGO_ENABLED=0 go build -o bin/cookiemonster-darwin-amd64 -v ./src/cmd

build-linux:
	GOOS=linux CGO_ENABLED=0 go build -o bin/cookiemonster-linux-amd64 -v ./src/cmd

clean:
	rm -rf ./bin ./vendor

docker:
	docker build --no-cache -t cookiemonster -f Dockerfile.cookiemonster .