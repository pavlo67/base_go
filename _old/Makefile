BUILD_TIME=$(shell date -u '+%Y-%m-%d %H:%M:%S')
CGO_ENABLED=0
COMMIT=$(shell git rev-parse --short HEAD)
TAG=$(shell git describe --tags $(git rev-list --tags --max-count=1))

LDFLAGS=-ldflags '-s -w -X "main.BuildTag=${TAG}" -X "main.BuildCommit=${COMMIT}" -X "main.BuildDate=${BUILD_TIME}"'

b:
	                        go build -o bin/demo           -v ${LDFLAGS} ./apps/demo/main

bl:
	GOOS=linux GOARCH=amd64 go build -o bin/demo           -v ${LDFLAGS} ./apps/demo/main


