# The Link Shortener

This project is a high-performance URL shortening service written in Go.

## Table of Contents
- [Requirements](#requirements)
- [Setup](#setup)
- [Running the Application](#running-the-application)
- [Makefile Commands](#makefile-commands)
- [Docker](#docker)
- [Postman](#postman)
- [Swagger](#swagger)

## Requirements

- Go 1.22+
- Docker & Docker Compose
- PostgreSQL

## Setup

1. Clone the repository:

```bash
git clone https://github.com
cd The-link-shortener-on-Golang
```

2. Create a `.env` file from the template:

```bash
make env
```

3. Configure the environment variables in the newly created `.env` file (e.g., database connection credentials and ports).

## Running the Application

### Quick Start via Docker
Deploy the entire infrastructure (application, PostgreSQL database, and Redis) with a single command:
```bash
docker-compose up -d
```

### Local Development Run
1. Start the database and redis infrastructure:
```bash
docker-compose up -d
```
2. Apply database migrations to create the required tables:
```bash
make migrate-up
```
3. Run the server locally:
```bash
make run
```

## Makefile Commands

The following commands are available via the `Makefile`:

##### `make build`: Build the project binary
##### `make run`: Run the server locally
##### `make test`: Run tests
##### `make clean`: Clean up binary files
##### `make lint`: Run `golangci-lint` static analysis tool
##### `make migrate-up`: Apply PostgreSQL database migrations
##### `make migrate-down`: Roll back PostgreSQL database migrations
##### `make env`: Create a `.env` file based on `.envexample`
##### `make help`: Display help information for available commands

## Docker

Use standard commands to manage your containers:

Start services in the background:
```bash
docker-compose up -d
```

Stop and remove containers:
```bash
docker-compose down
```

## Postman

Use the following payload example to send requests:

```json
{
  "url": "https://google.com",
  "custom_code": "my_custom_alias"
}
```

###### URL Shortener API Implementation:

- **POST Request** with the payload above to create a short link:
```bash
http://localhost:8080/shorten
```
*Note: The `custom_code` field is optional. If omitted, the server will automatically generate a random, unique 6-character code.*

- **GET Request** using the short code to redirect to the original website:
```bash
http://localhost:8080/my_custom_alias
```
*The request returns an HTTP `307 Temporary Redirect` status and instantly forwards the user.*

## Swagger

For visual testing and interactive API exploration, use Swagger UI:

```bash
http://localhost:8080/swagger/index.html
```

---

###### The server operates by creating a short link (with either a random or custom code) and saving it into a PostgreSQL database. Upon accessing the short link via a GET request, the application performs an instant HTTP redirect to the original URL and asynchronously pushes a tracking task into an internal Go channel (`clickQueue`). A pool of 5 active goroutine workers concurrently processes this queue to increment the click counter (`clicks`) for each link in the database, ensuring high throughput without blocking the main HTTP request-response cycle.
