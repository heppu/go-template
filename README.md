[![Go Reference](https://pkg.go.dev/badge/github.com/heppu/go-template.svg)](https://pkg.go.dev/github.com/heppu/go-template) [![codecov](https://codecov.io/github/heppu/go-template/graph/badge.svg?token=H3u7Ui9PfC)](https://codecov.io/github/heppu/go-template) ![main](https://github.com/heppu/go-template/actions/workflows/go.yaml/badge.svg?branch=main)

# Go project template

A comprehensive, production-ready Go project template enforcing best practices. This template provides a solid foundation for building scalable, observable, and maintainable Go services.

## Getting started

Fork this repo and create your project repo using the fork as a template. Clone the new repo and run `make -f rename.mk`. After that everyhting should be renamed based on your repository's name. To verify that renaming was successfull run `make all`. Create a PR with these changes and see how the github workflow gets triggered.

## Tooling

This template includes configuration for the following tools:

- **Go 1.24+**
- **Docker** - For containerization and local testing
- **Docker Compose** - For local testing
- **Make** - Build automation
- **golangci-lint** - Code quality and style enforcement
- **OpenTelemetry** - Distributed tracing and metrics
  - Jaeger UI: http://localhost:16686
  - Prometheus: http://localhost:9090
- **GitHub Actions** - CI workflows
- **Codecov** - Code coverage

## Libraries

Key libraries included and pre-configured:

- **ogen** - OpenAPI code generation with observability support
- **golang-migrate** - Database migration management
- **sqlx** - Mapping data between structs and SQL
- **testcontainers-go** - Integration testing with real services
- **stretchr/testify** - Test assertions

The template is structured to provide a solid foundation while allowing easy customization for your specific project needs.
