# dubplate version: v0.5.2

ROOT_DIR ?= $(shell git rev-parse --show-toplevel)
UNTRACKED ?= $(shell test -z "$(shell git ls-files --others --exclude-standard "$(ROOT_DIR)")" || echo -untracked)
VERSION ?= $(shell git describe --tags --dirty --always)$(UNTRACKED)

OS ?= linux
ARCH ?= amd64

BUILD_DIR ?= ./build/$(VERSION)/$(OS)-$(ARCH)

$(BUILD_DIR):
	mkdir -p $@

clean:
	rm $(BUILD_DIR)/*

version:
	@echo ${VERSION}

cmd-all: binary test-binary-version-output image

image:
	@echo skipping image build...
