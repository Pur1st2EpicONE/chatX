![chatX banner](assets/banner.png)

<h3 align="center">Simple CRUD service implemented in Go as an internship test project, with PostgreSQL storage and LRU caching for fast HTTP access.</h3>

<br>

## Table of Contents

- [Architecture](#architecture)  
- [Installation](#installation)  
- [Configuration](#configuration)  
- [Shutting down](#shutting-down)  
- [Request examples](#request-examples)  

<br>

## Architecture

- **App** — central orchestrator. Loads configuration, initializes logger, cache, storage, service, handlers and HTTP server, wires dependencies, and manages lifecycle and graceful shutdown via a shared context.

- **Handler (HTTP)** — Gin-based HTTP layer. Exposes REST endpoints under /api/v1/chats and serves Swagger UI at /swagger/\*any.

- **Service** — business logic layer. Validates input, enforces domain rules, coordinates cache and storage usage, and implements CRUD operations.

- **Cache** — in-memory LRU cache used to serve frequent reads with low latency. Configurable capacity and per-chat message limits.

- **Repository** — persistent data layer (PostgreSQL via GORM). Handles connection pooling and migrations (goose).

- **Logger** — structured JSON logger. Writes to log directory or stdout, supports debug/info/warn/error/fatal levels.

<br>

## Installation

⚠️ Note: This project requires Docker Compose, regardless of how you choose to run it.

First, clone the repository and enter the project folder:

```bash
git clone https://github.com/Pur1st2EpicONE/chatX.git
cd chatX
```

Then you have two options:

### 1. Run everything in containers

```bash
make
```

This will start the entire project fully containerized using Docker Compose.


### 2. Run chatX locally

```bash
make local
```

Starts PostgreSQL in container and runs the application locally. Useful for iterative development and debugging.

⚠️ Note: Local mode requires Go 1.25.1 installed on your machine.

<br>

## Configuration

### Runtime configuration

Service uses three configuration files, depending on the selected run mode:

[config.full.yaml](./configs/config.full.yaml) — used for the fully containerized setup

[config.dev.yaml](./configs/config.dev.yaml) — used for local development

[config.test.yaml](./configs/config.test.yaml) — used for testing

You may optionally review and adjust the corresponding configuration file to match your preferences. The default values are suitable for most use cases.

### Environment variables

Service uses a .env file for runtime configuration. You may create your own .env file manually before running the service, or edit [.env.example](.env.example) and let it be copied automatically on startup.
If environment file does not exist, .env.example is copied to create it. If environment file already exists, it is used as-is and will not be overwritten.

⚠️ Note: Keep .env.example for local runs. Some Makefile commands rely on it and may break if it's missing.

<br>

## Shutting down

Stopping chatX depends on how it was started:

- Local setup — press Ctrl+C to send SIGINT to the application. The service will gracefully close connections and finish any in-progress operations.  
- Full Docker setup — containers run by Docker Compose will be stopped automatically.

In both cases, to stop all services and clean up containers, run:

```bash
make down
```

⚠️ Note: In the full Docker setup, the log folder is created by the container as root and will not be removed automatically. To delete it manually, run:

```bash
sudo rm -rf <log-folder>
```

⚠️ Note: Docker Compose also creates a persistent volume for PostgreSQL data (chatx_postgres_data). This volume is not removed automatically when containers are stopped. To remove it and fully reset the environment, run:

```bash
make reset
```

<br>

## Request examples

⚠️ Note: You can explore and interact with the API via Swagger UI at host:port/swagger/index.html

### Create chat

```bash
curl -X POST http://localhost:8080/api/v1/chats/ \
  -H "Content-Type: application/json" \
  -d '{"title": "The best chat ever!!!"}'
```

Response:

```json
{
  "result": {
    "id": 1,
    "title": "The best chat ever!!!",
    "created_at": "2025-01-16T12:00:00Z"
  }
}
```

<br>

### Create message

```bash
curl -X POST http://localhost:8080/api/v1/chats/1/messages/ \
  -H "Content-Type: application/json" \
  -d '{"text": "Hi!"}'
```

Response:

```json
{
  "result": {
    "id": 10,
    "chat_id": 1,
    "text": "Hi!",
    "created_at": "2025-01-16T12:01:00Z"
  }
}
```

<br>

### Get chat

```bash
curl http://localhost:8080/api/v1/chats/1?limit=10
```

Response:

```json
{
  "result": {
    "id": 1,
    "title": "The best chat ever!!!",
    "created_at": "2025-01-16T12:00:00Z",
    "messages": [{ "id": 10, "chat_id": 1, "text": "Hi!", "created_at": "2025-01-16T12:01:00Z" }]
  }
}
```

<br>

### Delete chat

```bash
curl -X DELETE http://localhost:8080/api/v1/chats/1
```

Response:

```json
{ "result": "deleted" }
```
