# Build metadata
APPLICATION_NAME := "elcli"
APPLICATION_VERSION := "0.2.0"
APPLICATION_BUILDID=$(shell git rev-parse --short HEAD)

# Golang build parameters
GO_OS := ""
GO_ARCH := ""
GO_PKG := "github.com/unbrikd/edge-leap"
GO_BINDIR := "./bin"
GO_EXTENSION := ""

# Docker build parameters
DOCKER_FILE := "./docker/Dockerfile"
DOCKER_OPTS := --build-arg APPLICATION_VERSION=$(APPLICATION_VERSION) --build-arg APPLICATION_BUILDID=$(APPLICATION_BUILDID) --build-arg GO_PKG=$(GO_PKG) $(DOCKER_BUILD_EXTRAOPTS)
DOCKER_BUILD_EXTRAOPTS := ""
DOCKER_CONTEXT := "."

docker-image:
	@echo "---> Building docker image elcli:${APPLICATION_VERSION}"
	docker build $(DOCKER_OPTS) -t elcli:${APPLICATION_VERSION} -f $(DOCKER_FILE) .

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

show-version:
	@echo "v$(APPLICATION_VERSION), build $(APPLICATION_BUILDID)"

clean:
	@echo "---> Cleaning up"
	@echo "$(GO_BINDIR)/*" && rm -rf $(GO_BINDIR)/*