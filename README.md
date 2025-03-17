[![Go Reference](https://pkg.go.dev/badge/github.com/heppu/go-template.svg)](https://pkg.go.dev/github.com/heppu/go-template) [![codecov](https://codecov.io/github/heppu/go-template/graph/badge.svg?token=H3u7Ui9PfC)](https://codecov.io/github/heppu/go-template) ![main](https://github.com/heppu/go-template/actions/workflows/go.yaml/badge.svg?branch=main)

# Go project template

A comprehensive, production-ready Go project template enforcing best practices. This template provides a solid foundation for building scalable, observable, and maintainable Go services.

## Getting started

1. Fork this repo
2. Create project repo using the fork as a template
3. Clone the project repo and run `make -f rename.mk`
4. Run `make all` to verify that everything works
5. Create a PR see how the github workflow gets triggered

## Project layout

The template is structured to provide a solid foundation while allowing easy customization for your specific project needs.

### Go files

- `./api/` - rest api layer generated from `openapi.yaml`
- `./app/` - business logic that maps to rest endpoints inside api
- `./applicationtest/` - application/integration tests (tests executed against application binary)
- `./cmd/*/` - entry points for applications/binaries, each subdirectory becomes a separate executable
- `./server/` - configures http.Server with api handler
- `./store/` - database layer with migration support
- `./store/migrations/` - database schema migration files

### Non Go files

- `./.github/` - configuration files for github actions
- `./.golanci.yaml` - configuration golangci-linter
- `./.ogen.yaml` - configuration ogen api generator
- `./openapi.yaml` - api specification
- `./docker-compose.yaml` - configuration for services used in test
- `./Dockerfile` - image definition for Go binaries
- `./Makefile` - build tooling configuration
- `./ci.mk` - build tooling configuration for CI only targets
- `./rename.mk` - script to run rename after cloning initial template
- `./telemetry/` - configuration for otel related tools
- `./target/` - container for build and test artifacts

## Tooling

- **Make** - Build automation
- **Docker** - For containerization and local testing
- **Docker Compose** - For local testing
- **golangci-lint** - Code quality and style enforcement
- **gotestsum** - Test output formatter
- **GitHub Actions** - CI workflows
- **Codecov** - Code coverage

## Libraries

Key libraries included and pre-configured:

- **go-srvc/srvc** - Service library for life cycle management
- **go-srvc/mods** - Ready made modules for srvc
- **ogen** - OpenAPI code generation with observability support
- **golang-migrate** - Database migration management
- **sqlx** - Mapping data between structs and SQL
- **go-tstr/tstr** - Testing library with application test support
- **stretchr/testify** - Test assertions
