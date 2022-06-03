SHELL = bash
default: help

GIT_COMMIT := $(shell git rev-parse --short HEAD)
GIT_DIRTY := $(if $(shell git status --porcelain),+CHANGES)

GO_LDFLAGS := "-X github.com/hcjulz/damon/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY)"

HELP_FORMAT="    \033[36m%-25s\033[0m %s\n"
.PHONY: help
help: ## Display this usage information
	@echo "Valid targets:"
	@grep -E '^[^ ]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		sort | \
		awk 'BEGIN {FS = ":.*?## "}; \
			{printf $(HELP_FORMAT), $$1, $$2}'
	@echo ""

.PHONY: build
build:
	go build -o bin/damon ./cmd/damon

.PHONY: run
run:
	./bin/damon

.PHONY: install-osx
install-osx:
	cp ./bin/damon /usr/local/bin/damon

.PHONY: test
test:
	go test ./...

pkg/%/damon: GO_OUT ?= $@
pkg/windows_%/damon: GO_OUT = $@.exe
pkg/%/damon: ## Build Daemon for GOOS_GOARCH, e.g. pkg/linux_amd64/damon
	@echo "==> Building $@ with tags $(GO_TAGS)..."
	@CGO_ENABLED=0 \
		GOOS=$(firstword $(subst _, ,$*)) \
		GOARCH=$(lastword $(subst _, ,$*)) \
		go build -trimpath -ldflags $(GO_LDFLAGS) -tags "$(GO_TAGS)" -o $(GO_OUT) ./cmd/damon

.PRECIOUS: pkg/%/damon
pkg/%.zip: pkg/%/damon ## Build and zip Damon for GOOS_GOARCH, e.g. pkg/linux_amd64.zip
	@echo "==> Packaging for $@..."
	zip -j $@ $(dir $<)*

.PHONY: dev
dev: ## Build for the current development version
	@echo "==> Building damon..."
	@CGO_ENABLED=0 go build -ldflags $(GO_LDFLAGS) -o ./bin/damon ./cmd/damon
	@rm -f $(GOPATH)/bin/damon
	@cp ./bin/damon $(GOPATH)/bin/damon
	@echo "==> Done"

.PHONY: version
version:
ifneq (,$(wildcard version/version_ent.go))
	@$(CURDIR)/scripts/version.sh version/version.go version/version_ent.go
else
	@$(CURDIR)/scripts/version.sh version/version.go version/version.go
endif
