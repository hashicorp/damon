SHELL = bash

GIT_COMMIT=$$(git rev-parse --short HEAD)
GIT_BRANCH = $$(git branch --show-current)
GIT_SHA    = $$(git rev-parse HEAD)
GIT_DIRTY=$$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_IMPORT="github.com/hashicorp/nomad-pack/internal/pkg/version"
GO_LDFLAGS="-s -w -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"
VERSION = $(shell ./build-scripts/version.sh internal/pkg/version/version.go)

REPO_NAME    ?= $(shell basename "$(CURDIR)")
PRODUCT_NAME ?= $(REPO_NAME)
BIN_NAME     ?= $(PRODUCT_NAME)

# Get latest revision (no dirty check for now).
REVISION = $(shell git rev-parse HEAD)

# Get local ARCH; on Intel Mac, 'uname -m' returns x86_64 which we turn into amd64.
# Not using 'go env GOOS/GOARCH' here so 'make docker' will work without local Go install.
OS   = $(strip $(shell echo -n $${GOOS:-$$(uname | tr [[:upper:]] [[:lower:]])}))
ARCH = $(strip $(shell echo -n $${GOARCH:-$$(A=$$(uname -m); [ $$A = x86_64 ] && A=amd64 || [ $$A = aarch64 ] && A=arm64 ; echo $$A)}))
PLATFORM ?= $(OS)/$(ARCH)
DIST     = dist/$(PLATFORM)
BIN      = $(DIST)/$(BIN_NAME)

ifeq ($(firstword $(subst /, ,$(PLATFORM))), windows)
BIN = $(DIST)/$(BIN_NAME).exe
endif

build:
	go build -o bin/damon ./cmd/damon

run:
	./bin/damon

install-osx:
	cp ./bin/damon /usr/local/bin/damon

test:
	go test ./...


pkg/%/damon: GO_OUT ?= $@
pkg/%/damon: ## Build Nomad Autoscaler for GOOS_GOARCH, e.g. pkg/linux_amd64/nomad
	@echo "==> Building $@ with tags $(GO_TAGS)..."
	@CGO_ENABLED=0 \
		GOOS=$(firstword $(subst _, ,$*)) \
		GOARCH=$(lastword $(subst _, ,$*)) \
		go build -trimpath -ldflags $(GO_LDFLAGS) -tags "$(GO_TAGS)" -o $(GO_OUT)

pkg/windows_%/nomad-autoscaler: GO_OUT = $@.exe

# Common Dev make target.
# Includes GO_LDFLAGS for convenience.
# Deploys dev build to GOPATH.
.PHONY: dev
dev: GOPATH=$(shell go env GOPATH)
dev:
	@echo "==> Building damon..."
	@CGO_ENABLED=0 go build -ldflags $(GO_LDFLAGS) -o ./bin/damon ./cmd/damon
	@rm -f $(GOPATH)/bin/damon
	@cp ./bin/damon $(GOPATH)/bin/damon
	@echo "==> Done"

# Supports CRT expected directory structure
dist:
	mkdir -p $(DIST)
	echo '*' > dist/.gitignore

# Supports CRT builds
.PHONY: bin
bin: dist
	GOARCH=$(ARCH) GOOS=$(OS) go build -o $(BIN) ./cmd/damon

.PHONY: binpath
binpath:
	@echo -n "$(BIN)"

# Support CRT version inference
.PHONY: version
version:
	@echo $(VERSION)

# Supports Docker image builds.
export DOCKER_BUILDKIT=1
BUILD_ARGS = BIN_NAME=$(BIN_NAME) PRODUCT_VERSION=$(VERSION) PRODUCT_REVISION=$(REVISION)
TAG        = $(PRODUCT_NAME)/$(TARGET):$(VERSION)
BA_FLAGS   = $(addprefix --build-arg=,$(BUILD_ARGS))
FLAGS      = --target $(TARGET) --platform $(PLATFORM) --tag $(TAG) $(BA_FLAGS)

# Set OS to linux for all docker/* targets.
docker/%: OS = linux

# DOCKER_TARGET is a macro that generates the build and run make targets
# for a given Dockerfile target.
# Args: 1) Dockerfile target name (required).
#       2) Build prerequisites (optional).
define DOCKER_TARGET
.PHONY: docker/$(1)
docker/$(1): TARGET=$(1)
docker/$(1): $(2)
	docker build $$(FLAGS) .
	@echo 'Image built; run "docker run --rm $$(TAG)" to try it out.'

.PHONY: docker/$(1)/run
docker/$(1)/run: TARGET=$(1)
docker/$(1)/run: docker/$(1)
	docker run --rm $$(TAG)
endef

# Create docker/<target>[/run] targets.
$(eval $(call DOCKER_TARGET,dev,))
$(eval $(call DOCKER_TARGET,release,bin))

.PHONY: docker
docker: docker/dev

# Serveral of the following command rely on the following env vars being
# configured according to the CRT onboarding documentation.
# - ARTIFACTORY_TOKEN
# - ARTIFACTORY_USER
# - CRT_STAGING_REGISTRY
docker-login:
	echo "$(ARTIFACTORY_TOKEN)" | docker login -u $(ARTIFACTORY_USER) --password-stdin $(CRT_STAGING_REGISTRY)

docker-pull-staging:
	@docker pull $(CRT_STAGING_REGISTRY)/$(REPO_NAME)/release:$(VERSION)_$(GIT_SHA)

staging:
	@bob trigger-promotion \
	  --product-name=$(PRODUCT_NAME) \
	  --org=hashicorp \
	  --repo=$(REPO_NAME) \
	  --branch=$(GIT_BRANCH) \
	  --product-version=$(VERSION) \
	  --sha=$(GIT_SHA) \
	  --environment=nomad-oss \
	  --slack-channel=CUYKT2A73 \
	  staging

download:
	@bob download artifactory \
		-channel stable \
		-product-name $(PRODUCT_NAME) \
		-product-version $(VERSION) \
		-commit $(GIT_SHA)

verify-rpm:
	@docker run \
		-v $(CURDIR)/build-scripts:/local \
		-e ARTIFACTORY_TOKEN=$(ARTIFACTORY_TOKEN) \
		-e ARTIFACTORY_USER=$(ARTIFACTORY_USER) \
		-e REPO_NAME=$(REPO_NAME) \
		-e VERSION=$(VERSION) \
		-e GIT_SHA=$(GIT_SHA) \
		centos:7 \
		/bin/bash /local/verify-rpm.sh

debug-rpm:
	@docker run -it \
		-v $(CURDIR)/build-scripts:/local \
		-e ARTIFACTORY_TOKEN=$(ARTIFACTORY_TOKEN) \
		-e ARTIFACTORY_USER=$(ARTIFACTORY_USER) \
		-e REPO_NAME=$(REPO_NAME) \
		-e VERSION=$(VERSION) \
		-e GIT_SHA=$(GIT_SHA) \
		centos:7 \
		/bin/bash

verify-deb:
	@docker run \
		-v $(CURDIR)/build-scripts:/local \
		-e ARTIFACTORY_TOKEN=$(ARTIFACTORY_TOKEN) \
		-e ARTIFACTORY_USER=$(ARTIFACTORY_USER) \
		-e REPO_NAME=$(REPO_NAME) \
		-e VERSION=$(VERSION) \
		-e GIT_SHA=$(GIT_SHA) \
		ubuntu \
		/bin/bash /local/verify-deb.sh

debug-deb:
	@docker run -it \
		-v $(CURDIR)/build-scripts:/local \
		-e ARTIFACTORY_TOKEN=$(ARTIFACTORY_TOKEN) \
		-e ARTIFACTORY_USER=$(ARTIFACTORY_USER) \
		-e REPO_NAME=$(REPO_NAME) \
		-e VERSION=$(VERSION) \
		-e GIT_SHA=$(GIT_SHA) \
		ubuntu \
		/bin/bash
