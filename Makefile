# Build metadata
APPLICATION_NAME := "elcli"
APPLICATION_VERSION := "0.2.0"
APPLICATION_BUILDID=$(shell git rev-parse --short HEAD)
APPLICATION_ARCH := ""

# Golang build parameters
GO_OS := ""
GO_ARCH := ""
GO_PKG := "github.com/unbrikd/edge-leap"
GO_BINDIR := "./bin"
GO_EXTENSION := ""

# Docker build parameters
DOCKER_REPO := "ghcr.io/unbrikd/elcli"
DOCKER_FILE := "./docker/Dockerfile"
DOCKER_OPTS := --build-arg APPLICATION_VERSION=$(APPLICATION_VERSION) --build-arg APPLICATION_BUILDID=$(APPLICATION_BUILDID) --build-arg GO_PKG=$(GO_PKG) $(DOCKER_EXTRAOPTS)
DOCKER_EXTRAOPTS := ""
DOCKER_CONTEXT := "."

docker-prepare-buildx:
	docker buildx create \
	--name container-builder \
	--driver docker-container \
	--use --bootstrap

docker-buildx:
	docker buildx build \
    --platform linux/amd64,linux/arm64 \
    --push \
    --tag ghcr.io/unbrikd/elcli:latest \
    --tag ghcr.io/unbrikd/elcli:0.3.0 \
    -f ./docker/Dockerfile .

docker-image:
	@echo "---> Building docker image $(DOCKER_REPO):${APPLICATION_VERSION}-$(APPLICATION_BUILDID)$(APPLICATION_ARCH)"
	@docker build $(DOCKER_OPTS) -t $(DOCKER_REPO):${APPLICATION_VERSION}-$(APPLICATION_BUILDID)$(APPLICATION_ARCH) -f $(DOCKER_FILE) .

docker-linux:
	@echo "---> Building docker image for linux/amd64"
	@$(MAKE) docker-image DOCKER_EXTRAOPTS="--platform linux/amd64" APPLICATION_ARCH="-amd64"

	@echo "---> Building docker image for linux/arm64"
	@$(MAKE) docker-image DOCKER_EXTRAOPTS="--platform linux/arm64" APPLICATION_ARCH="-arm64"


build:
	@echo "---> $(GO_BINDIR)/elcli-v$(APPLICATION_VERSION).${GO_OS}-${GO_ARCH}$(EXTENSION)"
	@GOOS=$(GO_OS) GOARCH=$(GO_ARCH) go build \
		-ldflags "-s -w -X $(GO_PKG)/version.Version=$(APPLICATION_VERSION) -X $(GO_PKG)/version.Revision=$(APPLICATION_BUILDID)" \
		-o $(GO_BINDIR)/elcli-v$(APPLICATION_VERSION).${GO_OS}-${GO_ARCH}$(GO_EXTENSION)

build-macos:
	@echo "---> Building for darwin/amd64"
	@$(MAKE) build GO_OS=darwin GO_ARCH=amd64

	@echo "---> Building for darwin/arm64"
	@$(MAKE) build GO_OS=darwin GO_ARCH=arm64

build-linux:
	@echo "---> Building for linux/amd64"
	@$(MAKE) build GO_OS=linux GO_ARCH=amd64

	@echo "---> Building for linux/arm64"
	@$(MAKE) build GO_OS=linux GO_ARCH=arm64

build-windows:
	@echo "---> Building for windows/amd64"
	@$(MAKE) build GO_OS=windows GO_ARCH=amd64 GO_EXTENSION=".exe"

print-version:
	@echo $(APPLICATION_VERSION)

clean:
	@echo "---> Cleaning up"
	@echo "$(GO_BINDIR)/*" && rm -rf $(GO_BINDIR)/*
