# Build metadata
APPLICATION_NAME := "elcli"
APPLICATION_VERSION := "0.2.3"
APPLICATION_BUILDID=$(shell git rev-parse --short HEAD)
APPLICATION_ARCH := ""

# Golang build parameters
GO_OS := ""
GO_ARCH := ""
GO_MAIN := "./cmd/elcli.go"
GO_PKG := "github.com/unbrikd/edge-leap"
GO_BINDIR := "./bin"
GO_EXTENSION := ""

# Docker build parameters
DOCKER_REPO := "ghcr.io/unbrikd/elcli"
DOCKER_FILE := "./docker/Dockerfile"
DOCKER_OPTS := --build-arg APPLICATION_VERSION=$(APPLICATION_VERSION) --build-arg APPLICATION_BUILDID=$(APPLICATION_BUILDID) --build-arg GO_PKG=$(GO_PKG) --build-arg GO_MAIN=$(GO_MAIN) $(DOCKER_EXTRAOPTS)
DOCKER_EXTRAOPTS := ""
DOCKER_CONTEXT := "."

docker-buildx:
	@echo "---> Building docker image ghcr.io/unbrikd/$(APPLICATION_NAME):$(APPLICATION_VERSION)"
	@docker buildx build \
	--push \
    $(DOCKER_OPTS) \
    --tag ghcr.io/unbrikd/$(APPLICATION_NAME):latest \
    --tag ghcr.io/unbrikd/$(APPLICATION_NAME):$(APPLICATION_VERSION)\
    -f $(DOCKER_FILE) .

docker-buildx-allarch:
	@echo "---> Building docker image for linux/amd64 and linux/arm64"
	@$(MAKE) docker-buildx DOCKER_EXTRAOPTS="--platform linux/amd64,linux/arm64"

docker-image:
	@echo "---> Building docker image $(DOCKER_REPO):${APPLICATION_VERSION}-$(APPLICATION_BUILDID)$(APPLICATION_ARCH)"
	@docker build \
	$(DOCKER_OPTS) \
	-t $(DOCKER_REPO):${APPLICATION_VERSION}-$(APPLICATION_BUILDID)$(APPLICATION_ARCH) \
	-f $(DOCKER_FILE) .

build:
	@echo "---> $(GO_BINDIR)/elcli-v$(APPLICATION_VERSION).${GO_OS}-${GO_ARCH}$(EXTENSION)"
	@GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build \
		-ldflags "-s -w -X $(GO_PKG)/version.Version=$(APPLICATION_VERSION) -X $(GO_PKG)/version.Revision=$(APPLICATION_BUILDID)" \
		-o $(GO_BINDIR)/elcli-v$(APPLICATION_VERSION).${GO_OS}-${GO_ARCH}$(GO_EXTENSION) $(GO_MAIN)

build-macos:
	@echo "---> Building for darwin/amd64"
	@$(MAKE) build GO_OS=darwin GO_ARCH=amd64 GO_BINDIR=$(GO_BINDIR)

	@echo "---> Building for darwin/arm64"
	@$(MAKE) build GO_OS=darwin GO_ARCH=arm64 GO_BINDIR=$(GO_BINDIR)

build-linux:
	@echo "---> Building for linux/amd64"
	@$(MAKE) build GO_OS=linux GO_ARCH=amd64 GO_BINDIR=$(GO_BINDIR)

	@echo "---> Building for linux/arm64"
	@$(MAKE) build GO_OS=linux GO_ARCH=arm64 GO_BINDIR=$(GO_BINDIR)

build-windows:
	@echo "---> Building for windows/amd64"
	@$(MAKE) build GO_OS=windows GO_ARCH=amd64 GO_EXTENSION=".exe" GO_BINDIR=$(GO_BINDIR)

unit-tests:
	@echo "---> Running unit tests"
	@go test -v ./...

print-version:
	@echo $(APPLICATION_VERSION)

clean:
	@echo "---> Cleaning up"
	@echo "$(GO_BINDIR)/*" && rm -rf $(GO_BINDIR)/*
	@echo "./edge-leap.yaml" && rm -rf ./edge-leap.yaml
