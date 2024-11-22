SHELL = /bin/bash
GO := go
extra_env := $(GOENV)

export PKG := github.com/grokspawn/stencil
export GIT_COMMIT := $(or $(SOURCE_GIT_COMMIT),$(shell git rev-parse --short HEAD))
export STENCIL_VERSION := $(or $(SOURCE_GIT_TAG),$(shell git describe --always --tags HEAD))
export BUILD_DATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')



.PHONY: stencil
stencil:
	$(extra_env) $(GO) build $(extra_flags) $(TAGS) -o bin/stencil ./main.go
