VERSION ?= $(shell cat VERSION)
GIT_DIRTY = $(shell test -n "`git status --porcelain`" && echo "-dirty")

.PHONY: all
all: build

.PHONY: build
build: deps
	govvv build -pkg github.com/ryane/kfilt/cmd

.PHONY: docker
docker:
	docker build -t ryane/kfilt:${VERSION}${GIT_DIRTY} .

push: build
	docker push ryane/kfilt:${VERSION}${GIT_DIRTY}

.PHONY: deps
deps:
	GO111MODULE=off go get github.com/ahmetb/govvv
