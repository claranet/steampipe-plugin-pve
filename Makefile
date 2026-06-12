STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS = netgo

OS := $(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH := $(shell uname -m)
VERSION := $(shell cat VERSION)

ifeq ($(ARCH),x86_64)
	ARCH = amd64
endif
ifeq ($(ARCH),aarch64)
	ARCH = arm64
endif

PLUGIN_DIR = $(STEAMPIPE_INSTALL_DIR)/plugins/local/pve

.PHONY: build install clean test

PLATFORMS = linux-amd64 darwin-arm64
BUILD_TAGS ?=

define build_template
build_$(1): GOOS := $(word 1,$(subst -, ,$(1)))
build_$(1): GOARCH := $(word 2,$(subst -, ,$(1)))
build_$(1):
	go build -o steampipe-plugin-pve_v$(VERSION).$(1).plugin -tags "$(BUILD_TAGS)" *.go
endef

$(foreach platform,$(PLATFORMS),$(eval $(call build_template,$(platform))))

build: $(foreach platform,$(PLATFORMS),build_$(platform))

install: build
	mkdir -p $(PLUGIN_DIR)
	cp steampipe-plugin-pve.v$(VERSION).${OS}-${ARCH}.plugin $(PLUGIN_DIR)/steampipe-plugin-pve.plugin
	@echo ""
	@echo "Plugin installed to $(PLUGIN_DIR)"
	@echo ""
	@echo "Add this to ~/.steampipe/config/pve.spc:"
	@echo ""
	@echo '  connection "pve" {'
	@echo '    plugin = "local/pve"'
	@echo '    api_url = "https://your-pve:8006/api2/json"'
	@echo '    api_token_id = "user@realm!token"'
	@echo '    api_token_secret = "secret"'
	@echo '    tls_insecure = true'
	@echo '  }'


uninstall:
	rm -f $(PLUGIN_DIR)/steampipe-plugin-pve.plugin


clean:
	rm -f steampipe-plugin-pve*.plugin

test:
	go test ./... -v
