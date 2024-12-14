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

## Running the app

1. **Clone the repository:**

   ```bash
   git clone https://github.com/khaled4vokalz/gourl_shortener.git
   cd gourl_shortener
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Set up configuration:**
   Create a configuration file in YAML format or use environment variables for settings like the port, database, and cache.

4. **Run the application:**

   ```bash
   make run
   ```

   If you don't have `make` you can alternatively use

   ```bash
    go run ./cmd/gourl_shortener
   ```

## API Endpoints

### Shorten URL

**POST** `/shorten`

- without `Origin` header

  Request:

  ```bash
  curl --verbose --request POST --data '{"url": "http://example.com"}' localhost:8080/shorten
  ```

  Response:

  ```json
  { "shortened_url": "http://localhost:8080/CBXqmaO8" }
  ```

- with `Origin` header

  Request:

  ```bash
  curl --verbose --request POST --header "Origin: https://foo.com" --data '{"url": "http://example.com"}' localhost:8080/shorten
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

## Configuration

Configuration can be provided via YAML files or environment variables, currently it only supports config file in the `configuration` directory having the same name of the ENVIRONMENT env. Example YAML configuration:

```yaml
server:
  host: localhost
  port: 8080
db_conn_string: "user=postgres password=<pass> dbname=gourl_shortener sslmode=disable"
cache:
  host: localhost
  port: 6379
  database: 0
  username: <user>
  password: <pass>
shortener_props:
  length: 6 # the total bytes that should be considered from the SHA256 hash of the url
  max_attempt: 5 # maximum attempt the service should take when key collision happens for a url
```

## Testing

Run unit tests with:

```bash
make test
```

If you don't have `make` you can alternatively use

```bash
go test ./...
```

## Roadmap

- [] Implement analytics for shortened URLs (e.g., number of clicks)
- [] Add expiration time
- [] Add a web UI for managing shortened URLs
- [] Add support for Environment variables
- [] Add docker support
- [] House keeping
  - [] Use factories for db package
  - [] Provide default values for the configurations, probably my having a separate config manager

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
