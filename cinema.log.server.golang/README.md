

## How the project is structured


- cmd/api - contains the main.go
- internal - contains all code (each subfolder is a package)
- internal/database - change database.go if wanting to change connection to db e.g to postgres or another place
- internal/domain - all domain objects/structs e.g. users. This is used in a lot of places
- internal/server
- internal/server/routes.go - to set up new routes for the server
- internal/server/server.go - to set up dependency injection for all vertical slices (new db -> new store -> new service -> new handler)
- internal/{nameOfVerticalSlice} - each slice contains a handler, service and store with tests
