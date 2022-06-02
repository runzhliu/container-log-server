pkgs          = $(shell go list ./...)

PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)
DOCKER_IMAGE_NAME       ?= 	harbor.dev-fql.com/middleware/log-server

BRANCH      ?= $(shell git rev-parse --abbrev-ref HEAD)
BUILDDATE   ?= $(shell date --iso-8601=seconds)
BUILDUSER   ?= $(shell whoami)@$(shell hostname)
REVISION    ?= $(shell git rev-parse HEAD)
TAG_VERSION ?= $(shell git describe --tags --abbrev=0)

VERSION_LDFLAGS := \
  -X main.version=$(TAG_VERSION)

all: build

build:
	@echo ">> building code"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "$(VERSION_LDFLAGS)"

docker:
	@echo ">> building docker image"
	docker build -t "$(DOCKER_IMAGE_NAME):$(TAG_VERSION)" .
	docker push "$(DOCKER_IMAGE_NAME):$(TAG_VERSION)"

.PHONY: all build docker
