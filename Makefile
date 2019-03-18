SHELL := /bin/bash
REPO=hook
BINARY=lambda-fn

VERSION=0.1.0
BUILD_TIME=`date +%FT%T%z`

BRANCH=`git rev-parse --abbrev-ref HEAD`
COMMIT=`git rev-parse --short HEAD`

LDFLAGS="-X ${REPO}.version=${VERSION} -X ${REPO}.buildtime=${BUILD_TIME} -X ${REPO}.branch=${BRANCH} -X ${REPO}.commit=${COMMIT}"

default: build

.PHONY: clean
clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: pretest
pretest:
	@gofmt -d $$(find . -type f -name '*.go' -not -path "./vendor/*") 2>&1 | read; [ $$? == 1 ]

.PHONY: vet
vet:
	@go vet $(go list -f '{{ .ImportPath }}' ./... | grep -v vendor/)

.PHONY: test
test: pretest vet lint
	@go test -v $$(go list -f '{{ .ImportPath }}' ./... | grep -v vendor/) -p=1

.PHONY: fmt
fmt:
	@gofmt -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: lint
lint:
	@go get -v github.com/golang/lint/golint
	@golint ./... | grep -v vendor/ | true

.PHONY: build
build:
	@GOOS=linux go build -x -ldflags ${LDFLAGS} -o bin/${BINARY} github.com/umayr/${REPO}/cmd/${BINARY}
