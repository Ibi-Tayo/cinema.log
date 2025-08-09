# Project cinema.log.server.golang

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests

```bash
make all
```

Build the application

```bash
make build
```

Run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB Container

```bash
make docker-down
```

DB Integrations Test:

```bash
make itest
```

Live reload the application:

```bash
make watch
```

Run the test suite:

```bash
make test
```

Clean up binary from the last build:

```bash
make clean
```

## How the project is structured

cmd/api - contains the main.go
internal - contains all code (each subfolder is a package)
internal/database - change database.go if wanting to change connection to db e.g to postgres or another place
internal/domain - all domain objects/structs e.g. users. This is used in a lot of places
internal/server
internal/server/routes.go - to set up new routes for the server
internal/server/server.go - to set up dependency injection for all vertical slices (new db -> new store -> new service -> new handler)
internal/{nameOfVerticalSlice} - each slice contains a handler, service and store
