# dubplate version: v0.5.2

OUTBIN ?= $(BUILD_DIR)/$(APP_NAME)

VERSION_VAR ?= main.version
LDFLAGS = -ldflags "-w -X $(VERSION_VAR)=$(VERSION)"
GOBUILD_FLAGS ?= -installsuffix cgo $(LDFLAGS) -o $(OUTBIN)
GOBUILD_ENVVARS ?= CGO_ENABLED=0 GOOS=$(OS) GOARCH=$(ARCH)
GOBUILD_CMD ?= $(GOBUILD_ENVVARS) go build $(GOBUILD_FLAGS)

dummy:
	@echo No default rule set yet

binary: $(BUILD_DIR)
	$(GOBUILD_CMD) ./cmd/$(APP_NAME)

binaries: $(BINARIES)

$(BINARIES):
	$(MAKE) cmd-all \
		APP_NAME=$@

test-binary-version-output: VERSION_CMD ?= $(OUTBIN) version
test-binary-version-output:
	@echo testing output of $(VERSION_CMD)
	test "$(shell $(VERSION_CMD))" = "$(VERSION)" && echo PASSED
