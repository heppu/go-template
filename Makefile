## Requirements:
#	- go
#	- docker
#	- awk
#	- printf
#	- cut
#	- uniq

### Environment variables

# For scratch image to work CGO has to be disabled.
# If you need to use CGO, you can override this and use a different base image.
export CGO_ENABLED ?= 0

### Static variables
DIST_DIR          := target/dist
TEST_DIR          := target/test
UNIT_COV_DIR      := ${TEST_DIR}/unit
UNIT_BIN_COV_DIR  := ${UNIT_COV_DIR}/cov/bin
UNIT_TXT_COV_DIR  := ${UNIT_COV_DIR}/cov/txt
UNIT_JUNIT_DIR    := ${UNIT_COV_DIR}/junit
APP_COV_DIR       := ${TEST_DIR}/application
APP_BIN_DIR       := ${APP_COV_DIR}/cov/bin
APP_TXT_COV_DIR   := ${APP_COV_DIR}/cov/txt
APP_JUNIT_DIR     := ${APP_COV_DIR}/junit
CMB_COV_DIR       := ${TEST_DIR}/combined
CMB_TXT_COV_DIR   := ${CMB_COV_DIR}/cov/txt
CMD_DIR           := ./cmd
BUILD_TARGETS     := $(wildcard ${CMD_DIR}/*)
TARGET_NAMES      := $(BUILD_TARGETS:${CMD_DIR}/%=%)
BIN_TARGETS       := ${TARGET_NAMES:%=bin/%}
IMG_TARGETS       := ${TARGET_NAMES:%=img/%}
PUSH_TARGETS      := ${TARGET_NAMES:%=.push/%}
NOOP              :=
SPACE             := ${NOOP} ${NOOP}

### Build variables
TARGET            = ${@F}
TARGET_DIR        = ${DIST_DIR}/${TARGET}
TARGET_BIN        = ${TARGET_DIR}/${TARGET}
TARGET_PKG        = ${CMD_DIR}/${TARGET}

### Override these in CI
IMG_REG           ?=
IMG_REPO          ?=
IMG_NAME		  ?= $(subst ${SPACE},/,$(filter-out ,$(strip ${IMG_REG} ${IMG_REPO} ${TARGET})))
IMG_TAGS          ?= dev

### Docker build variables
IMG_TARGET_ARGS = ${IMG_TAGS:%=-t ${IMG_NAME}:%}
IMG_BUILD_ARGS  = --build-arg TARGET=${TARGET}

foo:
	echo ${IMG_NAME}

.DEFAULT_GOAL := help
.PHONY: help
help: ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nStatic targets:\n"} /^[a-zA-Z0-9_\/-]+:.*?##/ { printf "  \033[36m%-20s\033[0m  %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@for target in ${BIN_TARGETS}; do printf "  \033[36m%-22s\033[0m %s\n" $${target} "Build binary for $${target}"; done
	@for target in ${IMG_TARGETS}; do printf "  \033[36m%-22s\033[0m %s\n" $${target} "Build image for $${target}"; done

.PHONY: all
all: clean generate lint test bin img ## Test and build all targets

.PHONY: telemetry-up
telemetry-up:
	docker compose -f ./telemetry/docker-compose.yaml up -d
	@printf "\nJaeger UI available at http://localhost:16686\n"
	@printf "Grafana UI available at http://localhost:16686\n"

.PHONY: telemetry-down
telemetry-down:
	docker compose -f ./telemetry/docker-compose.yaml down

.PHONY: clean
clean: ## Clean up environment
	rm -rf ./target

.PHONY: generate
generate: ## Run code generators
	go generate ./...

.PHONY: lint
lint: ## Run linter
	go tool github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

.PHONY: test
test: test/unit test/application ## Run all tests
	rm -rf ${CMB_TXT_COV_DIR}
	mkdir -p ${CMB_TXT_COV_DIR}
	go tool covdata textfmt -i=${APP_BIN_DIR},${UNIT_BIN_COV_DIR} -o ${CMB_TXT_COV_DIR}/cover.txt
	go tool cover -html=${CMB_TXT_COV_DIR}/cover.txt

.PHONY: test/unit
test/unit: ## Run unit tests
	rm -rf ${UNIT_BIN_COV_DIR} ${UNIT_TXT_COV_DIR} ${UNIT_JUNIT_DIR}
	mkdir -p ${UNIT_BIN_COV_DIR} ${UNIT_TXT_COV_DIR} ${UNIT_JUNIT_DIR}
	CGO_ENABLED=1 go tool gotestsum --junitfile=${UNIT_JUNIT_DIR}/junit.xml -- -race -covermode=atomic -coverprofile=${UNIT_TXT_COV_DIR}/cover.txt ./... -test.gocoverdir=$(abspath ${UNIT_BIN_COV_DIR})

.PHONY: test/application
test/application: ## Run application tests
	rm -rf ${APP_BIN_DIR} ${APP_TXT_COV_DIR} ${APP_JUNIT_DIR}
	mkdir -p ${APP_BIN_DIR} ${APP_TXT_COV_DIR} ${APP_JUNIT_DIR}
	GOCOVERDIR=$(abspath ${APP_BIN_DIR}) go tool gotestsum --junitfile=${APP_JUNIT_DIR}/junit.xml -- -tags=applicationtest -count=1 ./applicationtest/...
	go tool covdata textfmt -i=${APP_BIN_DIR} -o ${APP_TXT_COV_DIR}/cover.txt

.PHONY: bin
bin: ${BIN_TARGETS} ## Build all binaries

.PHONY: img
img: ${IMG_TARGETS} ## Build all images


### Dynamic targets
#
# Each package under ./cmd/ will have a corresponding target to build a binary and an image.

.PHONY: ${BIN_TARGETS}
${BIN_TARGETS}:
	mkdir -p ${TARGET_DIR}
	go build -o  ${TARGET_BIN} ${TARGET_PKG}

.PHONY: ${IMG_TARGETS}
${IMG_TARGETS}:
	docker buildx build -f ./Dockerfile ${IMG_BUILD_ARGS} ${IMG_TARGET_ARGS} ${TARGET_DIR}

.PHONY: ${PUSH_TARGETS}
${PUSH_TARGETS}:
	docker push --all-tags ${IMG_NAME}

### Targets meant only for CI
#
# These targets use .{target} notation so that they don't show up in autocomplete.
.PHONY: .push
.push: ${PUSH_TARGETS} ## Push all images

.PHONY: .pr-check
.pr-check: .check-duplicated-migrations .check-modified-migrations .check-generated-code

.PHONY: .main-check
.main-check: .check-duplicated-migrations .check-generated-code

.PHONY: .check-duplicated-migrations
.ONESHELL:
SHELL = /bin/bash
.check-duplicated-migrations:
	@echo "Performing duplicated migration check"
	output="$$(ls -1 store/migrations/ | cut -d "_" -f1 | uniq -D)"
	if [[ -n $$output ]]; then
		echo "Found duplicate migration versions:"
		echo "$$output"
		exit 1
	fi
	echo "No duplicated migrations found"

.PHONY: .check-modified-migrations
.ONESHELL:
SHELL = /bin/bash
.check-modified-migrations:
	@if test -z "$$BASE_REF"; then
		echo "BASE_REF must be set"
		exit 1
	fi
	if test -z "$$HEAD_REF"; then
		echo "HEAD_REF must be set"
		exit 1
	fi
	echo "Performing migration verification on PR against $$BASE_REF"
	git fetch origin $$BASE_REF
	git fetch origin $$HEAD_REF
	git diff --exit-code --name-only --diff-filter=D origin/$$BASE_REF origin/$$HEAD_REF -- store/migrations/ || (echo "main branch has new migrations, please rebase" && exit 1)
	echo "No modified migrations found"

.PHONY: .check-generated-code
.check-generated-code: generate
	@git diff --exit-code --name-only || (printf "Running make generate modified code base.\nRun make genrate before commiting" && exit 1)
