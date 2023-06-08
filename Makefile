APP=$(shell basename $(shell git remote get-url origin))
DOCKERREGISTRY=ghcr.io/ivanloktionov
VERSION=v$(shell git describe --tags --abbrev=0)-$(shell git rev-parse --short HEAD)
TARGETOS=linux
TARGETARCH=arm64

format:
	gofmt -s -w ./

lint:
	golint

test:
	go test -v

get:
	go get


linux:
	GOARCH=amd64 GOOS=linux go build -o kbot-linux main.go 
	
arm:
	GOARCH=amd64 GOOS=arm go build -o kbot-arm main.go

macos:
	GOARCH=amd64 GOOS=macos go build -o kbot-macos main.go

windows:
	GOARCH=amd64 GOOS=windows go build -o kbot-windows main.go

build: format
    CGO_ENABLED=0 GOARCH=${shell dpkg --print-architecture} go build -v -o kbot -ldflags "-X="github.com/ivanloktionov/kbot/cmd.appVersion=${VERSION}

image: 
	docker build . -t ${DOCKERREGISTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH} 

push: 
	docker push ${DOCKERREGISTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH}

clean:
	rm -rf kbot

	docker image ls | grep -v REPOSITORY| tr -s ' '|cut -f 3 -d ' '| head -1
