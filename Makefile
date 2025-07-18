export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

# https://stackoverflow.com/questions/4122831/disable-make-builtin-rules-and-variables-from-inside-the-make-file
MAKEFLAGS += --no-builtin-rules
.SUFFIXES:
APP_NAME			  := "your app name"
BUILD_DATE            := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT            := $(shell git rev-parse HEAD)
GIT_REMOTE            := origin
GIT_BRANCH            := $(shell git rev-parse --symbolic-full-name --verify --quiet --abbrev-ref HEAD)
GIT_TAG               := $(shell git describe --tags --abbrev=0  2> /dev/null || echo untagged)
GIT_TREE_STATE        := $(shell if [ -z "`git status --porcelain`" ]; then echo "clean"; else echo "dirty"; fi)
RELEASE_TAG           := $(shell if [[ "$(GIT_TAG)" =~ ^v[0-9]+\.[0-9]+\.[0-9]+.*$$ ]]; then echo "true"; else echo "false"; fi)
DEV_BRANCH            := $(shell [ $(GIT_BRANCH) = master ] || [ `echo $(GIT_BRANCH) | cut -c -8` = release- ] || [ `echo $(GIT_BRANCH) | cut -c -4` = dev- ] || [ $(RELEASE_TAG) = true ] && echo false || echo true)
SRC                   := $(pwd)

GREP_LOGS             := ""

VERSION               := 0.0.1
DOCKER_PUSH           := true

# VERSION is the version to be used for files in manifests and should always be latest unless we are releasing
# we assume HEAD means you are on a tag
#ifeq ($(RELEASE_TAG),true)
#VERSION               := $(GIT_TAG)
#endif

CGO_ENABLED ?= 1
WASM_ENABLED ?= 1

GO := CGO_ENABLED=$(CGO_ENABLED) GO111MODULE=on go
GOVERSION := $(shell cat ./.go-version)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
DISABLE_CGO := CGO_ENABLED=0


IMAGE := "your image repos"/$(APP_NAME)

override LDFLAGS = -X myapi/pkg/util.version=$(VERSION) \
-X myapi/pkg/util.buildDate=${BUILD_DATE} \
-X myapi/pkg/util.gitCommit=${GIT_COMMIT} \
-X myapi/pkg/util.gitTreeState=${GIT_TREE_STATE}

ifneq ($(GIT_TAG),)
override LDFLAGS += -X myapi/pkg/util.gitTag=${GIT_TAG}
endif

ifndef $(GOPATH)
GOPATH=$(shell go env GOPATH)
export GOPATH
endif

#.PHONY: cli
#cli: ./main.go
#	go build -v -ldflags '${LDFLAGS}' -o $(APP_NAME)

.PHONY: image-static
image-static:
CGO_ENABLED=0 WASM_ENABLED=0 $(MAKE) build-linux-static
@$(MAKE) image-quick-static

.PHONY: image-quick-static
image-quick-static:
sed -e 's/GOARCH/amd64/g' Dockerfile > .Dockerfile_amd64
docker build -t $(IMAGE):$(VERSION) -f .Dockerfile_amd64 .

build:
$(GO) build $(GO_TAGS) -o $(APP_NAME)_$(GOOS)_$(GOARCH) -ldflags '$(LDFLAGS)' main.go

.PHONY: build-linux-static
build-linux-static:
@$(MAKE) GOOS=linux GOARCH=amd64 build  WASM_ENABLED=0 CGO_ENABLED=0