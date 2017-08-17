PHONY: all clean clean-build clean-dist clean-tmp dockes-shim dte tf-apply tf-plan
.DEFAULT_GOAL := help

help:
	$(info available targets:)
	@awk '/^[a-zA-Z\-\_0-9\.\$$\(\)\%/]+:/ { \
		helpMsg = $$0; \
		nb = sub(/^[^:]*:.* ## /, "", helpMsg); \
		if (nb) \
			print  $$1 "\t" helpMsg; \
	}' \
	$(MAKEFILE_LIST) | column -ts $$'\t' | \
	grep --color '^[^ ]*'

# project variables
PROJ_NAME := pr0nbot
VERSION = $(if $(SNAP),SNAPSHOT,$(shell git describe --always --dirty))

# helper variables
BUILD := build
DIST := dist
LAMBDA_DIR := lambdas
TF := terraform
TMP := $(BUILD)/tmp
REPO_PATH := github.com/theherk/pr0nbot
LDFLAGS = "-X $(REPO_PATH)/util.Version=$(VERSION)"

LAMBDA_ARCS :=
define LAMBDA_template
LAMBDA_ARCS += $(DIST)/$(1).zip
endef

LAMBDAS := test-lambda
$(foreach LAMBDA,$(LAMBDAS),$(eval $(call LAMBDA_template,$(LAMBDA))))

ifeq ($(OS),Windows_NT)
	MKDIR := mkdir
	FixPath = $(subst /,\,$1)
else
	MKDIR := mkdir -p
	FixPath = $1
endif

t:
	echo $(OS)
	$(MKDIR) $(call FixPath,test/path)

all: bin lambda-arcs ## build binary and lambda distributions

bin: $(BUILD)/$(PROJ_NAME) ## build the main program / cli

clean: clean-build clean-dist clean-tmp ## remove build, dist, and tmp directories

clean-build: ## remove build directory
	rm -rf $(BUILD)

clean-dist: ## remove build directory
	rm -rf $(DIST)

clean-tmp: ## remove temporary directory
	rm -rf $(TMP)

docker-shim: ## pull docker image for building lambdas
	docker pull eawsy/aws-lambda-go-shim:latest

dte: $(TF)/api.tf ## generate api.tf from downtoearth.json

lambda-arcs: $(LAMBDA_ARCS) ## build all plugins and distribution archives

test: test-unit test-integration ## run all tests

test-unit: ## unit tests
	go test ./cmd/... ./lib/... -v -tags unit

test-integration: ## integration tests
	@echo "No integration tests to run"

tf-apply: ## apply environment state changes
	terraform apply -refresh=true -state=$(TF)/terraform.tfstate $(TF)/

tf-plan: ## show environment state changes
	terraform plan -refresh=true -state=$(TF)/terraform.tfstate $(TF)/

$(BUILD)/$(PROJ_NAME): ## build the main program / cli
	$(MKDIR) $(call FixPath,$(TMP))
	go build -i -v -ldflags=$(LDFLAGS) -o $@

$(DIST)/%.zip: clean-tmp docker-shim ## lambda distribution archive
	$(MKDIR) $(call FixPath,$(BUILD)/$*)
	$(MKDIR) $(call FixPath, $(DIST))
	$(MKDIR) $(call FixPath, $(TMP)/handler)
	docker run --rm \
	  -v $(GOPATH):/go \
	  -v $(CURDIR):/work \
	  -w /work \
	  eawsy/aws-lambda-go-shim:latest make $(TMP)/$*.zip
	mv $(TMP)/$*.zip $@

$(TF)/api.tf: ## generate api.tf from downtoearth.json
	downtoearth generate terraform/downtoearth.json -c terraform/$(PROJ_NAME)-root.tf

# WARNING: The following targets are expected to be run inside docker.
# DO NOT run directly.

.PRECIOUS: $(BUILD)/%/handler.so
$(BUILD)/%/handler.so:
	go build -buildmode=plugin -ldflags='-w -s' -o $@ $(LAMBDA_DIR)/$*/handler.go
	@chown $(shell stat -c '%u:%g' .) $@

$(TMP)/%.zip: $(BUILD)/%/handler.so
	@cp $< $(TMP)/
	@cp /shim/__init__.pyc $(TMP)/handler/__init__.pyc
	@cp /shim/proxy.pyc $(TMP)/handler/proxy.pyc
	@cp /shim/runtime.so $(TMP)/handler/runtime.so
	@find $(TMP)/ -exec touch -t 201302210800 {} +
	cd $(TMP) && zip -qrX $(notdir $@) * ; cd $(CURDIR)
	@chown -R $(shell stat -c '%u:%g' .) $(TMP)
