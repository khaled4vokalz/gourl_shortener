## Table of Contents

- [URL Shortener in Go](#url-shortener-in-go)
  - [Features](#features)
  - [Prerequisites](#prerequisites)
  - [Running the app](#running-the-app)
    - [Backend](#backend)
    - [Frontend](#frontend)
    - [Using docker-compose](#using-docker-compose)
  - [API Endpoints](#api-endpoints)
    - [Shorten URL](#shorten-url)
    - [Get Original URL](#get-original-url)
    - [Backend Health Check](#backend-health-check)
  - [Configuration](#configuration)
    - [Backend](#backend-1)
    - [Frontend](#frontend-1)
  - [Testing](#testing)
    - [Backend](#backend-2)
  - [Roadmap](#roadmap)
  - [License](#license)

# URL Shortener in Go

This project is a URL Shortener service implemented in Go. It provides RESTful APIs to shorten long URLs and retrieve original URLs using a shortened identifier. The project uses an extensible design with support for multiple storage backends, caching, and configuration management. This is a Go 🐹 learning project for me

## Features

- **Shorten URLs**: Generate a short, unique URL for any valid original URL.
- **Retrieve Original URLs**: Retrieve (get re-directed to) the original URL using the shortened identifier.
- **In-Memory Storage**: Default storage backend for rapid prototyping.
- **PostgreSQL Support**: Optional storage backend for persistence.
- **Caching with Redis**: Reduce database loads by caching
- **Environment Configurations**: Support for development and production environments.
- **Extensible Design**: Easily swap or add new storage backends.

## Prerequisites

- [Go](https://golang.org/dl/) 1.20+
- [Redis](https://redis.io/) for caching (optional, but recommended)
- [PostgreSQL](https://www.postgresql.org/) for persistent storage (optional)
- [Make](https://www.gnu.org/software/make/manual/make.html) - optional
- [Docker](https://www.docker.com/) - optional

## Running the app

**Clone the repository:**

```bash
$ git clone https://github.com/khaled4vokalz/gourl_shortener.git
```

### Backend

- **Bare Metal:**

  - **get into server directory**

    ```bash
    $ cd gourl_shortener/server
    ```

  - **install dependencies:**

    ```bash
    $ go mod tidy
    ```

  - **set up configuration:**
    create a configuration file in yaml format or use environment variables for settings like the port, database, and cache.
  - If you're using postgres as storage option, then run the `scripts/init-db-script.sh` script, so that needed schemas are created. Make sure to pass in the `DB_USER`, `DB_PASSWORD` and `DB_NAME` envs.

  - **run the application:**

    ```bash
    $ make run
    ```

    if you don't have `make` you can alternatively use

    ```bash
    $ go run ./cmd/gourl_shortener
    ```

- **Docker**

  - build the docker image

    ```bash
    $ docker build --tag go-url-shortener-server .
    ```

  - run the container

    ```bash
    $ docker run --detach --env GOURLAPP_storage_type=in-memory --env GOURLAPP_cache_enabled=false --name gourl_shortener --publish 8082:8080 go-url-shortener-server
    ```

### Frontend

- **Bare Metal:**

  - Get into the client directory

    ```bash
    $ cd client
    ```

  - Install dependencies (use node 20+)

    ```bash
    $ nvm use 20.0.0
    $ npm ci
    ```

  - Start the app

    We can set `REACT_APP_BACKEND_URL` to the url where the backend api is available, e.g. `REACT_APP_BACKEND_URL=http://localhost:8080` before running `npm start`. By default the app uses `http://localhost:8082` as the backend url.

    ```bash
    $ npm start
    ```

  App should be running in development mode. Open [http://localhost:3000](http://localhost:3000) to view it in the browser. It's using Hot reload.

- **Docker**

  - Get into the client directory

    ```bash
    $ cd client
    ```

  - build the docker image

    Build time arg `BACKEND_URL` may be set to override the default `http://localhost:8082` that the app uses as backend url.

    ```bash
    $ docker build --build-arg BACKEND_URL=http://url-where-backend-is-running:port --tag go-url-shortener-client .
    ```

  - run the container

    ```bash
    $ docker run --detach --name gourl_shortener_client --publish 3000:80 go-url-shortener-client
    ```

  Open [http://localhost:3000](http://localhost:3000) to view it in the browser.

### Using docker-compose

Running the below command should spin up the containers with needed database and tables in it.

- The postgres db is listening at port `5433`
- The redis db is listening at port `6380`
- The backend is listening at port `8082`
- The frontend app is listening at `3000`

```bash
$ POSTGRES_PASSWORD=<your-postgres-pass> DB_PASSWORD=<your-db-pass> docker compose up --detach
```

## API Endpoints

### Shorten URL

**POST** `/shorten`

- without `Origin` header

  Request:

  ```bash
  curl --verbose --request POST --data '{"url": "http://example.com"}' localhost:8082/shorten
  ```

  Response:

  ```json
  { "shortened_url": "http://localhost:8082/CBXqmaO8" }
  ```

- with `Origin` header

  Request:

  ```bash
  curl --verbose --request POST --header "Origin: https://foo.com" --data '{"url": "http://example.com"}' localhost:8082/shorten
  ```

  Response:

  ```json
  { "shortened_url": "https://foo.com/CBXqmaO8" }
  ```

### Get Original URL

**GET** `/CBXqmaO8`

Response headers:

```sh
< HTTP/1.1 308 Permanent Redirect
< Content-Type: text/html; charset=utf-8
< Location: http://example.com
```

### Backend Health Check

**GET** `/health/live`

Response (with status code `200`):

```json
{ "databaes_is_live": true, "cache_is_alive": true }
```

Response (with status code `500`, in case either/all of the components are not alive):

```json
{ "databaes_is_live": true, "cache_is_alive": false }
```

## Configuration

### Backend

- Configuration can be provided via YAML files or environment variables, currently it only supports config file in the `configuration` directory having the same name of the ENVIRONMENT env. Example YAML configuration:

```yaml
server:
  host: localhost
  port: 8080
storage:
  type: postgres
  dbConnString: "user=shortener password=<pass> dbname=gourl_shortener sslmode=disable"
cache:
  enabled: true
  host: localhost
  port: 6379
  database: 0
shortenerProps:
  length: 6 # the total bytes that should be considered from the SHA256 hash of the url
  maxAttempt: 5 # maximum attempt the service should take when key collision happens for a url
environment: dev # the environment of the application, either `dev` or `prod`.
```

- Override configurations using Environment Variables

Any configuration mentioned above can be overridden using Environment variables using `GOURLAPP_` prefix and use the property tree separated by underscores (`_`). e.g. If we want to override `maxAttempt` configuration, we can set it like `GOURLAPP_shortenerProps_maxAttempt=10`

### Frontend

- only `REACT_APP_BACKEND_URL` environment variable that can override the backend url the app should talk to

## Testing

### Backend

Get into server directory

```bash
$ cd gourl_shortener/server
```

Run unit tests with:

```bash
$ make test
```

If you don't have `make` you can alternatively use

```bash
$ go test ./...
```

## Roadmap

- [] Implement analytics for shortened URLs (e.g., number of clicks)
- [x] Add expiration time
- [x] Add a web UI for managing shortened URLs
- [x] Add support for Environment variables in the config files
- [x] Add docker support
- [x] Add logging
- [x] Health-check
- [] House keeping
  - [] Provide default values for the configurations, probably my having a separate config manager

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
