# Project related variables
PKG = 1000eyes_exporter
APP_NAME=1000eyes_exporter

# Directories
WD := $(subst $(BSLASH),$(FSLASH),$(shell pwd))
MD := $(subst $(BSLASH),$(FSLASH),$(shell dirname "$(realpath $(lastword $(MAKEFILE_LIST)))"))
BUILD_DIR = $(WD)/build
PKG_DIR = $(MD)
CMD_DIR = $(PKG_DIR)/cmd
DIST_DIR = $(WD)/dist
LOG_DIR = $(WD)/log
REPORT_DIR = $(WD)/reports

# Make Variables
M = $(shell printf "\033[34;1m▶\033[0m")
DONE="$(M) done ✨"
VERSION := $(shell git describe --exact-match --tags 2>/dev/null)
ifndef VERSION
	VERSION := dev
endif
GIT_TAG := $(shell git describe --exact-match --tags 2>git_describe_error.tmp; rm -f git_describe_error.tmp)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_COMMIT := $(shell git rev-parse HEAD)

GOBIN = $(shell go env GOPATH)/bin
ARCHES ?= amd64
OSES ?= linux darwin
OUTTPL = $(BUILD_DIR)/$(APP_NAME)-$(VERSION)-{{.OS}}_{{.Arch}}
GZCMD = tar -czf
ZIPCMD = zip
SHACMD = sha256sum
VET_RPT=vet.out
COVERAGE_RPT=coverage.out

LDFLAGS = -X $(PKG)/version.APP_NAME=$(APP_NAME) \
	-X $(PKG)/version.commit=$(GIT_COMMIT) \
	-X $(PKG)/version.branch=$(GIT_BRANCH) \
	-X $(PKG)/version.version=$(VERSION) \
	-X $(PKG)/version.buildTime=$(shell date -Iseconds)

## deps: Download and Install any missing dependecies
.PHONY: deps
deps:
	go mod download
	@echo $(DONE) "-- Deps"

## tidy: Verifies and downloads all required dependencies
.PHONY: tidy
tidy:
	@echo "$(M) 🏃 go mod tidy..."
	@mkdir -pv $(REPORT_DIR)
	go mod verify
	go mod tidy
	@if ! git diff --quiet; then \
		echo "WARNING:  'go mod tidy' resulted in changes or working tree is dirty. See diff.out for details"; \
		git --no-pager diff > $(REPORT_DIR)/diff.out; \
	fi
	@echo $(DONE) "-- Tidy"

## fmt: Runs gofmt on all source files
.PHONY: fmt
fmt:
	@echo "$(M) 🏃 gofmt..."
	@ret=0 && for d in $$(go list -f '{{.Dir}}' ./...); do \
		gofmt -l -w $$d/*.go || ret=$$? ; \
	 done ; exit $$ret
	@echo $(DONE) "-- Fmt"

## clean: Removes build, dist and report dirs
.PHONY: clean
clean:
	@echo "$(M)  🧹 Cleaning build ..."
	go clean $(PKG) || true
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -rf $(REPORT_DIR)
	@echo $(DONE) "-- Clean"

## build: Install deps, format then builds binary (linux/amd64, darwin/amd64) in ./build
.PHONY: build
build: deps tidy fmt
	@mkdir -pv $(BUILD_DIR)
	@echo "$(LDFLAGS)"
	@echo "  $(M)  Checking if there is any missing dependencies...\n"
	@$(MAKE) deps
	@echo "  $(M)  Building...\n"
	@echo "GOBIN: $(GOBIN)"
	$(GOBIN)/gox -arch="$(ARCHES)" -os="$(OSES)" -output="$(OUTTPL)/{{.Dir}}" \
      	-tags "$(BUILD_TAGS)" -ldflags "$(LDFLAGS)"
	$(info "Built version:$(VERSION), build:$(GIT_COMMIT)")
	@echo $(DONE) "-- Build"

## debug: Print make env information
.PHONY: debug
debug:
	$(info PKG=$(PKG))
	$(info APP_NAME=$(APP_NAME))
	$(info MD=$(MD))
	$(info WD=$(WD))
	$(info PKG_DIR=$(PKG_DIR))
	$(info CMD_DIR=$(CMD_DIR))
	$(info BUILD_DIR=$(BUILD_DIR))
	$(info DIST_DIR=$(DIST_DIR))
	$(info LOG_DIR=$(LOG_DIR))
	$(info REPORT_DIR=$(REPORT_DIR))
	$(info VET_RPT=$(VET_RPT))
	$(info COVERAGE_RPT=$(COVERAGE_RPT))
	$(info VERSION=$(VERSION))
	$(info GIT_COMMIT=$(GIT_COMMIT))
	$(info GIT_TAG=$(GIT_TAG))
	$(info GOBIN=$(GOBIN))
	$(info ARCHES=$(ARCHES))
	$(info OSES=$(OSES))
	$(info LDFLAGS=$(LDFLAGS))
	@echo $(DONE) "-- Debug"

.PHONY: help
help: Makefile
	@echo "\n Choose a command run in "$(PROJECTNAME)":\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'