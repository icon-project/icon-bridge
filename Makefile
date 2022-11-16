#-------------------------------------------------------------------------------
#
# 	Makefile for building target binaries.
#

# Configuration
BUILD_ROOT = $(abspath ./)
BIN_DIR = ./bin
LINUX_BIN_DIR = ./build/linux

GOBUILD = go build
GOBUILD_TAGS =
# GOBUILD_ENVS = CGO_ENABLED=0
GOBUILD_LDFLAGS =
GOBUILD_FLAGS = -tags "$(GOBUILD_TAGS)" -ldflags "$(GOBUILD_LDFLAGS)"
GOBUILD_ENVS_LINUX = $(GOBUILD_ENVS) GOOS=linux GOARCH=amd64

GOTEST = go test
GOTEST_FLAGS = -test.short

# Build flags
GL_VERSION ?= $(shell git describe --always --tags --dirty)
GL_TAG ?= latest
BUILD_INFO = $(shell go env GOOS)/$(shell go env GOARCH) tags($(GOBUILD_TAGS))-$(shell date '+%Y-%m-%d-%H:%M:%S')

#
# Build scripts for command binaries.
#
CMDS = $(patsubst cmd/%,%,$(wildcard cmd/*))
.PHONY: $(CMDS) $(addsuffix -linux,$(CMDS))
define CMD_template
$(BIN_DIR)/$(1) $(1) : GOBUILD_LDFLAGS+=$$($(1)_LDFLAGS)
$(BIN_DIR)/$(1) $(1) :
	@ \
	rm -f $(BIN_DIR)/$(1) ; \
	echo "[#] go build ./cmd/$(1)"
	$$(GOBUILD_ENVS) \
	go build $$(GOBUILD_FLAGS) \
	    -o $(BIN_DIR)/$(1) ./cmd/$(1)

$(LINUX_BIN_DIR)/$(1) $(1)-linux : GOBUILD_LDFLAGS+=$$($(1)_LDFLAGS)
$(LINUX_BIN_DIR)/$(1) $(1)-linux :
	@ \
	rm -f $(LINUX_BIN_DIR)/$(1) ; \
	echo "[#] go build ./cmd/$(1)"
	$$(GOBUILD_ENVS_LINUX) \
	go build $$(GOBUILD_FLAGS) \
	    -o $(LINUX_BIN_DIR)/$(1) ./cmd/$(1)
endef
$(foreach M,$(CMDS),$(eval $(call CMD_template,$(M))))

# Build flags for each command
iconbridge_LDFLAGS = -X 'main.version=$(GL_VERSION)' -X 'main.build=$(BUILD_INFO)'
BUILD_TARGETS += iconbridge

linux : $(addsuffix -linux,$(BUILD_TARGETS))

SOLIDITY_DIST_DIR = $(BUILD_ROOT)/build/contracts/solidity

$(SOLIDITY_DIST_DIR)/%:
	$(eval MODULE := $(patsubst $(SOLIDITY_DIST_DIR)/%,%,$@))
	mkdir -p $@/contracts ; \
	mkdir -p $@/migrations ; \
	cp -r solidity/$(MODULE)/contracts/* $@/contracts ; \
	cp -r solidity/$(MODULE)/migrations/* $@/migrations ; \
	cp solidity/$(MODULE)/{truffle-config.js,*.json,*.lock} $@/ ; \
	cp solidity/$(MODULE)/truffle-config.js $@/ ; \
	cp solidity/$(MODULE)/*.json $@/ ; \
	cp solidity/$(MODULE)/*.lock $@/ ; \
	rm -rf $@/contracts/Mock ; \

dist-sol-bmc: $(SOLIDITY_DIST_DIR)/bmc
	cd $(SOLIDITY_DIST_DIR)/bmc ; 

dist-sol-bmv: $(SOLIDITY_DIST_DIR)/bmv
	cd $(SOLIDITY_DIST_DIR)/bmv ; 

dist-sol-bts: $(SOLIDITY_DIST_DIR)/bts
	cd $(SOLIDITY_DIST_DIR)/bts ;

# dist-sol-native_bsh: $(SOLIDITY_DIST_DIR)/bsh
# 	cd $(SOLIDITY_DIST_DIR)/bsh ;

dist-sol-token_bsh: $(SOLIDITY_DIST_DIR)/TokenBSH
	cd $(SOLIDITY_DIST_DIR)/TokenBSH ;

clean-dist-sol:
	rm -rf $(SOLIDITY_DIST_DIR)

dist-sol: dist-sol-bmc dist-sol-bts

dist-javascore:
	$(BUILD_ROOT)/docker/javascore/build.sh

ICONBRIDGE_IMAGE = iconbridge:$(GL_TAG)
ICONBRIDGE_DOCKER_DIR = $(BUILD_ROOT)/build/iconbridge

iconbridge-image: iconbridge-linux
	@ echo "[#] Building image $(ICONBRIDGE_IMAGE) for $(GL_VERSION)"
	@ rm -rf $(ICONBRIDGE_DOCKER_DIR)
	@ \
	BIN_DIR=$(abspath $(LINUX_BIN_DIR)) \
	BIN_VERSION=$(GL_VERSION) \
	BUILD_TAGS="$(GOBUILD_TAGS)" \
	DIST_DIR="$(BUILD_ROOT)/build/contracts" \
	$(BUILD_ROOT)/docker/iconbridge/build.sh $(ICONBRIDGE_IMAGE) $(BUILD_ROOT) $(ICONBRIDGE_DOCKER_DIR)

.PHONY: test

test :
	$(GOBUILD_ENVS) $(GOTEST) $(GOBUILD_FLAGS) ./... $(GOTEST_FLAGS)

.DEFAULT_GOAL := all
all : $(BUILD_TARGETS)
