# Dreams API

This repository contains the backend API for the [dreams-ui](https://github.com/bbrighter/dreams-ui) frontend.

The API is designed to record and organize dreams. Its main focus is on managing dream descriptions, categories, and the people associated with them. It also provides basic statistics to help analyze the recorded dreams.

## Installation and local setup

The project requires Go 1.27 or later.

Clone the repository and download its dependencies, then install all files via:

```bash
go mod download
```

Start the API locally with:

```bash
go run .
```

The application runs on `localhost:5000` by default. Database migrations are applied automatically when the application starts.

## Docker

The repository includes a [Dockerfile](/Dockerfile) for building the API as a lightweight Alpine-based image.

Build the image locally with:

```bash
docker build -t dreams-api .
```

New Docker images are published automatically to [Docker Hub](https://hub.docker.com/r/bbrighter/dreams-api) when a Git tag starting with `v` is pushed to GitHub. For example:

The publication workflow is defined in [.github/workflows/ci.yml](/.github/workflows/ci.yml). It runs the tests and license checks first, then builds and publishes the `linux/arm64` image as `bbrighter/dreams-api`.

## Configuration

Configuration is loaded from environment variables in [config/config.go](/config/config.go). The available settings and their defaults are:

| Environment variable | Default | Description |
| --- | --- | --- |
| `HOST_IP_ADDRESS` | `localhost` | Host address on which the API listens |
| `HOST_PORT` | `5000` | Port on which the API listens |
| `LOG_LEVEL` | `info` | Application log level |
| `DB_NAME` | `dreams.sqlite` | Database name inside the folder `data` |


For example, to run the API on all local interfaces and a different port:

```bash
HOST_IP_ADDRESS=0.0.0.0 HOST_PORT=8080 go run .
```

## Tests

Run the test suite with:

```bash
go test -v ./...
```

## Lefthook

This repository uses [Lefthook](https://github.com/evilmartians/lefthook) for Git hooks. The `pre-push` hook runs the test suite and checks that dependencies use an approved license.

After installing Lefthook, register the hooks in the local repository with:

```bash
go tool lefthook install
```


## OpenAPI documentation

The generated OpenAPI documentation is available in the [docs](/docs/) directory:
To regenerate the OpenAPI specification, install the Swag CLI if necessary:


Then run the VS Code task named `Create OpenAPI Spec`. This task executes:

```bash
go tool swag init --requiredByDefault
```

# Production run

Easiest run it in production with docker, e.g. via

```dockerfile
networks:
  network:

volumes:
  caddy_data:
  caddy_config:

services:

  dreams-api:
    image: bbrighter/dreams-api:v7.1.0
    restart: unless-stopped
    environment:
      - GIN_MODE=release
      - HOST_IP_ADDRESS=dreams-api
      - HOST_PORT=5000
    networks:
      - network
    volumes:
      - ./dreams.sqlite:/app/data/dreams.sqlite

  dreams-ui:
    image: bbrighter/dreams-ui:v7.0.1
    restart: unless-stopped
    networks:
      - network
    volumes:
      - ./dreamsConfig.json:/usr/share/nginx/html/config.json


  caddy:
   image: caddy:2.10.0-alpine
   container_name: apps-caddy
   restart: unless-stopped
   ports:
     - "80:80"
     - "443:443"
   volumes:
     - ./Caddyfile:/etc/caddy/Caddyfile
     - caddy_data:/data
     - caddy_config:/config
     - ./certs:/etc/caddy/certs:ro
   environment:
     - CADDY_HOST=host.docker.internal
   networks:
     - network
   depends_on:
     - dreams-ui
     - dreams-api
```