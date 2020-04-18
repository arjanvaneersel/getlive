# GetLive

Copyright 2020, GiveLive team  
info@example.com

## Licensing

TODO

## Project structure

The project follows the common directory structure for Go projects.

In `cmd` are subdirectories for all executable commands and services.
In `internal` are all internal packages.
In `internal/platform` are all internal packages, which might eventually need their own repo for code reusability.

## How to launch a development environment

1. Clone the repo
2. Run `make all` to build the givelive-api and metrics containers.
3. Run `make up` to launch the dev environment consisting, consisting of the getlive api service, tracing (zipkin) and metrics.
4. Run `make seed` to run migration and seeding of test data.
   
### Test environment endpoints

- API: http://localhost:3000/v1
- Debug: http://localhost:4000
- Zipkin tracing: http://localhost:9411
- Metrics: 
  - Expvar: http://localhost:3001
  - Debug: http://localhost:4001

### Test environment users
`admin@example.com` and `user@example.com`. The password for both users is `gophers`.

### Other actions

- `make down` to stop the dev environment
- `make test` to run all unit and integration tests
- `make migrate` to run only the migrations (without seeding)
- `make metrics` to (re)build only the metrics container
- `make getlive-api` to (re)build only the API container

## API endpoints

For now take a look at `cmd/getlive-api/internal/handlers/routes.go`.

## TODO
- [ ] Document API endpoints
- [ ] Improve security: Only owner or admin can delete entries
- [ ] Improve security: Only admin can approve entries
- [ ] Aggregator service: Web scraper
- [ ] Temporary admin interface (web or CLI) for manual approval of entries until there's a real frontend