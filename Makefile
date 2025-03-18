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
NOOP              :=
SPACE             := ${NOOP} ${NOOP}

### Build variables
TARGET            = demo
TARGET_DIR        = ${DIST_DIR}/${TARGET}
TARGET_BIN        = ${TARGET_DIR}/${TARGET}
TARGET_PKG        = ./cmd/${TARGET}

### Override these in CI
IMG_REG           ?=
IMG_REPO          ?=
IMG_NAME		  ?= $(subst ${SPACE},/,$(filter-out ,$(strip ${IMG_REG} ${IMG_REPO} ${TARGET})))
IMG_TAGS          ?= dev

### Docker build variables
IMG_TARGET_ARGS = ${IMG_TAGS:%=-t ${IMG_NAME}:%}
IMG_BUILD_ARGS  = --build-arg TARGET=${TARGET}

### OpenTelemetry variables
OTEL_ENV_VARS := OTEL_EXPORTER_OTLP_ENDPOINT="http://127.0.0.1:4318" OTEL_EXPORTER_OTLP_PROTO=http OTEL_EXPORTER_OTLP_INSECURE=true OTEL_SERVICE_NAME=${TARGET}

### Swagger UI
SWAGGER_UI_VERSION := 5.20.1
SWAGGER_UI_DIR     := ./server/swaggerui
SWAGGER_OLD_URL    := https://petstore.swagger.io/v2/swagger.json
SWAGGER_NEW_URL    := /docs/openapi.yaml

.DEFAULT_GOAL := help
.PHONY: help
help: ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nStatic targets:\n"} /^[a-zA-Z0-9_\/-]+:.*?##/ { printf "  \033[36m%-20s\033[0m  %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: all
all: clean .WAIT generate .WAIT lint test bin .WAIT img ## Test and build all targets

.PHONY: bin
bin: ## Build binary
	mkdir -p ${TARGET_DIR}
	go build -o  ${TARGET_BIN} ${TARGET_PKG}

.PHONY: img
img: ## Build image
	docker buildx build -f ./Dockerfile ${IMG_BUILD_ARGS} ${IMG_TARGET_ARGS} ${TARGET_DIR}

.PHONY: run
run: telemetry-up db-up ## Run application
	@printf "Starting server at http://127.0.0.1:8080\n"
	@printf "Swagger UI available at http://127.0.0.1:8080/docs/swaggerui\n"
	${OTEL_ENV_VARS} \
	API_ADDR=127.0.0.1:8080 \
	go run ${TARGET_PKG}

.PHONY: test
test: test/unit test/app ## Run all tests
	rm -rf ${CMB_TXT_COV_DIR}
	mkdir -p ${CMB_TXT_COV_DIR}
	go tool covdata textfmt -i=${APP_BIN_DIR},${UNIT_BIN_COV_DIR} -o ${CMB_TXT_COV_DIR}/cover.txt
	go tool cover -html=${CMB_TXT_COV_DIR}/cover.txt

.PHONY: test/unit
test/unit: ## Run unit tests
	rm -rf ${UNIT_BIN_COV_DIR} ${UNIT_TXT_COV_DIR} ${UNIT_JUNIT_DIR}
	mkdir -p ${UNIT_BIN_COV_DIR} ${UNIT_TXT_COV_DIR} ${UNIT_JUNIT_DIR}
	CGO_ENABLED=1 go tool gotestsum --junitfile=${UNIT_JUNIT_DIR}/junit.xml -- -race -covermode=atomic -coverprofile=${UNIT_TXT_COV_DIR}/cover.txt ./... -test.gocoverdir=$(abspath ${UNIT_BIN_COV_DIR})

.PHONY: test/app
test/app: ## Run application tests
	rm -rf ${APP_BIN_DIR} ${APP_TXT_COV_DIR} ${APP_JUNIT_DIR}
	mkdir -p ${APP_BIN_DIR} ${APP_TXT_COV_DIR} ${APP_JUNIT_DIR}
	GOCOVERDIR=$(abspath ${APP_BIN_DIR}) go tool gotestsum --junitfile=${APP_JUNIT_DIR}/junit.xml -- -tags=applicationtest -count=1 ./applicationtest/...
	go tool covdata textfmt -i=${APP_BIN_DIR} -o ${APP_TXT_COV_DIR}/cover.txt

.PHONY: test/app-otel
test/app-otel: telemetry-up ## Run application tests with OpenTelemetry
	${OTEL_ENV_VARS} $(MAKE) test/app

.PHONY: lint
lint: ## Run linter
	go tool github.com/golangci/golangci-lint/cmd/golangci-lint run ./...

.PHONY: telemetry-up
telemetry-up: ## Start telemetry stack
	docker compose -f ./telemetry/docker-compose.yaml up -d --wait
	@printf "\nJaeger UI available at http://127.0.0.1:16686\n"

.PHONY: telemetry-down
telemetry-down: ## Stop telemetry stack
	docker compose -f ./telemetry/docker-compose.yaml down

.PHONY: db-up
db-up: ## Start db
	docker compose up -d --wait

.PHONY: db-down
db-down: ## Stop db
	docker compose down

.PHONY: generate
generate: ## Run code generators
	go generate ./...

update-swagger-ui: ## Update Swagger UI
	rm -rf ${SWAGGER_UI_DIR}
	mkdir -p ${SWAGGER_UI_DIR}
	curl -s -L https://github.com/swagger-api/swagger-ui/archive/refs/tags/v${SWAGGER_UI_VERSION}.tar.gz | \
		tar -zxv --strip-components=2 -C ${SWAGGER_UI_DIR} swagger-ui-${SWAGGER_UI_VERSION}/dist/
	rm ${SWAGGER_UI_DIR}/*.map
	sed -i 's|${SWAGGER_OLD_URL}|${SWAGGER_NEW_URL}|g' ./server/swaggerui/swagger-initializer.js

.PHONY: clean
clean: ## Clean up environment
	rm -rf ./target
