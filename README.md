# Config Center

A Service to get config by application version and platform

## Install

Copy this repository on your local computer:

`git clone git@github.com:ashandi/config_center.git`

## Local development

Use the following command to run the service locally from the root dir of the project:

`make local-run`

Docker is required. This command will raise up two containers - the one with service itself, the second - with the database.

The database will be already initialized with some data for tests.

The local service will be available on `localhost:8080`

## Usage

`GET /config` to get config

See for the more details in `api/swagger.yaml`

## Add new dependencies

1. Describe new dependency in `api/swagger.yaml` and update the `types.ConfigResponse`.
2. Create new `RequestHandler` in `internal/server/handlers/config_request_handler`. It should implement interface RequestHandler.
Here you can add your validation rules for the request parameters and add new dependency
to the `Response`.
3. Add your `RequestHandler` to the chain on `BaseHandler` initialization.

## Tests

To run project tests, use the following command from the root dir of the project:

`go test ./...`