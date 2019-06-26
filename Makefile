GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "-dirty")
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

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
	CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags "-w -X github.com/ryane/kfilt/cmd.Version=${GIT_BRANCH} -X github.com/ryane/kfilt/cmd.GitCommit=${GIT_SHA}${GIT_DIRTY}" .
	docker build -t ryane/kfilt:${GIT_SHA}${GIT_DIRTY} .

push: build
	docker push ryane/kfilt:${GIT_SHA}${GIT_DIRTY}
