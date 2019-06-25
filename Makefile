GIT_SHA   = $(shell git rev-parse --short HEAD)
GIT_DIRTY = $(shell test -n "`git status --porcelain`" && echo "-dirty")

.PHONY: all
all: test build

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test -v ./...

.PHONY: docker
docker:
	docker build -t ryane/kfilt:${GIT_SHA}${GIT_DIRTY} .

push: build
	docker push ryane/kfilt:${GIT_SHA}${GIT_DIRTY}
