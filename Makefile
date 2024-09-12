# Build metadata
GIT_COMMIT=$(shell git rev-parse --short HEAD)
APPLICATION_VERSION := "0.2.0"

# Docker build variables
DOCKERFILE := "./docker/Dockerfile"

# Golang build variables
GOOS := ""
GOARCH := ""
GO_PKG := "github.com/unbrikd/edge-leap"
DESTDIR := "./bin"

docker-image:
	@echo "---> Building docker image for $(GOOS)/$(GOARCH)"
	@docker build --platform=$(GOOS)/$(GOARCH) -t elcli -f $(DOCKERFILE) .

build:
	@echo "---> $(DESTDIR)/elcli-v$(APPLICATION_VERSION).${GOOS}-${GOARCH}$(EXTENSION)"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-s -w -X $(GO_PKG)/version.Version=$(APPLICATION_VERSION) -X $(GO_PKG)/version.Revision=$(GIT_COMMIT)" \
		-o $(DESTDIR)/elcli-v$(APPLICATION_VERSION).${GOOS}-${GOARCH}$(EXTENSION)

build-macos:
	@echo "---> Building for darwin/amd64"
	@$(MAKE) build GOOS=darwin GOARCH=amd64

	@echo "---> Building for darwin/arm64"
	@$(MAKE) build GOOS=darwin GOARCH=arm64

build-linux:
	@echo "---> Building for linux/amd64"
	@$(MAKE) build GOOS=linux GOARCH=amd64

	@echo "---> Building for linux/arm64"
	@$(MAKE) build GOOS=linux GOARCH=arm64

build-windows:
	@echo "---> Building for windows/amd64"
	@$(MAKE) build GOOS=windows GOARCH=amd64 EXTENSION=".exe"

print-version:
	@echo $(APPLICATION_VERSION)

clean:
	@echo "---> Cleaning up"
	@rm -rf $(DESTDIR)/*