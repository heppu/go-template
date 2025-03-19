[![Go Reference](https://pkg.go.dev/badge/github.com/heppu/go-template.svg)](https://pkg.go.dev/github.com/heppu/go-template) [![codecov](https://codecov.io/github/heppu/go-template/graph/badge.svg?token=H3u7Ui9PfC)](https://codecov.io/github/heppu/go-template) ![main](https://github.com/heppu/go-template/actions/workflows/go.yaml/badge.svg?branch=main)

# Go project template

This template provides a solid foundation for building scalable, observable, and maintainable Go services.

## Getting started

### Creating new project

1. Fork this repo
2. Create project repo using the fork as a template
3. Clone the project repo and run `make -f rename.mk`
4. Run `make all` to verify that everything works
5. Create a PR see how the github workflow gets triggered

### Adding endpoints

1. Add paths and components in [server/openapi.yaml](server/openapi.yaml)
2. Run `make generate`
3. Implement missing methods in [app/app.go]

### Adding tables

1. Add new migration file `{NUM}_{NAME}.up.sql` eg `002_add_users_table.up.sql` under [store/migrations](store/migrations)
2. Write `CREATE TABLE` statement in that file
3. Now your migration is automatically applied when application starts

### Example

You can test this with simple example spec by running:

```sh
mv ./server/openapi-simple.yaml ./server/openapi.yaml
make generate
```

And now just implement missing methods in `./app/app.go`.

## Project layout

The template is structured to provide a solid foundation while allowing easy customization for your specific project needs.

### Go files

- `./api/` - rest api layer generated from [openapi.yaml](server/openapi.yaml)
- `./app/` - business logic that maps to rest endpoints inside api
- `./applicationtest/` - application/integration tests (tests executed against application binary)
- `./cmd/demo/` - main package for application, automatically renamed based on repo name
- `./server/` - configures http.Server with api handler
- `./store/` - database layer with migration support
- `./store/migrations/` - database schema migration files

### Non Go files

- `./.github/` - github actions workflow files
- `./.golanci.yaml` - golangci-linter configuration
- `./.ogen.yaml` - ogen generator configuration
- `./.codecov.yml` - codecov configuration
- `./server/openapi.yaml` - api specification
- `./server/swaggerui` - swagger ui files
- `./docker-compose.yaml` - configuration for services used in test
- `./Dockerfile` - image definition for Go binaries
- `./Makefile` - build tooling configuration
- `./ci.mk` - build tooling configuration for CI only targets
- `./rename.mk` - script to run rename after cloning initial template
- `./telemetry/` - configuration for otel related tools
- `./target/` - container for build and test artifacts

## Tooling

### Libraries

- [go-srvc/srvc](https://github.com/go-srvc/srvc) - Service library for life cycle management
- [go-srvc/mods](https://github.com/go-srvc/mods) - Ready made modules for srvc
- [golang-migrate](https://github.com/go-tstr/tstr/) - Database migration management
- [jmoiron/sqlx](https://github.com/jmoiron/sqlx) - Mapping data between structs and SQL
- [go-tstr/tstr](https://github.com/go-tstr/tstr) - Testing library with application test support
- [stretchr/testify](https://github.com/stretchr/testify) - Test assertions
- [automaxprocs](github.com/uber-go/automaxprocs) - Set runtime CPU resources automatically
- [automemlimit](github.com/KimMachineGun/automemlimit) - Set runtime MEM resources automatically

### External Tools

- [Make](https://www.gnu.org/software/make/) - Build automation
- [Docker](https://docs.docker.com/engine/install/) - For containerization and local testing
- [Docker Compose](https://docs.docker.com/compose/install/) - For local testing

### Go Tools (using [go tool](https://go.dev/doc/modules/managing-dependencies#tools))

- [golangci-lint](https://golangci-lint.run/) - Code quality and style enforcement
- [gotestsum](https://github.com/gotestyourself/gotestsum) - Test output formatter
- [ogen](https://ogen.dev/docs/intro) - OpenAPI code generation with observability and validation

### CI

- [GitHub Actions](https://docs.github.com/en/actions) - CI workflows
- [Codecov](https://app.codecov.io/github/heppu/go-template) - Code coverage
