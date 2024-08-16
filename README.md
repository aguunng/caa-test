# Caa Test Backend Service

## Overview

Caa Test service. You can read detailed information about this service [here](/docs/README.md)

### Environment Variables

To run this project, you will need to add the following environment variables to your `.env` file

```
QISCUS_APP_ID=
QISCUS_SECRET_KEY=
QISCUS_OMNICHANNEL_URL=
```

### Run Locally

To run the project locally, follow these steps:

- Clone this repository
- Navigate to the directory
- Format code and tidy modfile: `make tidy`
- Run test: `make test`, make sure that all tests are passing
- Run the server: `make run bin=server`, or run the application with reloading on file changes with: `make run/live bin=server`. You can also apply this to the cron application by changing the parameter to `bin=cron`
- The backend server will be accessible at `http://localhost:8080`
- You can find another useful commands in `Makefile`
- Alternatif run `go run main.go api`

### Generate Mock for Service

- Install [Mockery](https://github.com/vektra/mockery)
- Add the following code in the service file: `//go:generate mockery --all --case snake --output ./mocks --exported`
- Run go generate using `make generate`

### URL Demo

For a demo link of this project [here](https://caa-test.fly.dev/)