# dubplate version: v0.15.0

OUTBIN ?= $(BUILD_DIR)/$(APP_NAME)

VERSION_VAR ?= main.version
LDFLAGS = -ldflags "-w -X $(VERSION_VAR)=$(VERSION)"
GOBUILD_FLAGS ?= -installsuffix cgo $(LDFLAGS) -o $(OUTBIN)
GOBUILD_ENVVARS ?= CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH)
GOBUILD_CMD ?= $(GOBUILD_ENVVARS) go build $(GOBUILD_FLAGS)

# go install cannot cross-compile when GOBIN is set, so install builds for the host
GOINSTALL_ENVVARS ?= CGO_ENABLED=0
GOINSTALL_FLAGS ?= $(LDFLAGS)
GOINSTALL_CMD ?= $(GOINSTALL_ENVVARS) go install $(GOINSTALL_FLAGS)

dummy:
	@echo No default rule set yet

binary: $(BUILD_DIR)
	$(GOBUILD_CMD) ./cmd/$(APP_NAME)

binaries: $(COMPONENTS:=-binary)

$(COMPONENTS:=-binary):
	$(MAKE) binary \
		APP_NAME=$(@:-binary=)

.PHONY: install installs

install:
	$(GOINSTALL_CMD) ./cmd/$(APP_NAME)

installs: $(COMPONENTS:=-install)

$(COMPONENTS:=-install):
	$(MAKE) install \
		APP_NAME=$(@:-install=)


test-binary-version-output: VERSION_CMD ?= $(OUTBIN) version
test-binary-version-output:
	@echo testing output of $(VERSION_CMD)
	test "$(shell $(VERSION_CMD))" = "$(VERSION)" && echo PASSED
