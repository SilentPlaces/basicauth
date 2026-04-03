# Basic Authentication in Golang

This project provides a JWT-based authentication API built with Go, using Gin as the HTTP framework and a hexagonal-style architecture.

## Features

- User registration with email verification flow.
- User login and JWT token generation.
- Refresh token endpoint for rotating access tokens.
- Protected user endpoint with JWT middleware.
- Externalized configuration through Consul and secrets via Vault.
- MySQL persistence and Redis for verification token state.

## Architecture

The codebase is organized in a hexagonal-style structure:

- `cmd`: Application entry point.
- `internal/application/port`: Use-case dependency contracts (ports).
- `internal/application/usecase`: Core application use-cases.
- `internal/adapters/inbound/http/gin`: HTTP handlers, middleware, and router (Gin).
- `internal/infrastructure/di`: Dependency composition root.
- `internal/repositories`: Data access implementations.
- `internal/services`: Service implementations (auth, registration, mail, consul, vault, etc.).
- `internal/db`: MySQL and Redis client initialization.
- `migrations`: Database schema migrations.

## API Routes

- `POST /auth/login`
- `POST /auth/refresh-token`
- `GET /user` (requires `Authorization: Bearer <token>`)
- `POST /register/init`
- `POST /register/verify`
- `POST /register/resend-verification`

## Main Dependencies

- `gin-gonic/gin`: HTTP framework and router.
- `golang-jwt/jwt`: JWT creation and validation.
- `go-sql-driver/mysql`: MySQL driver.
- `redis/go-redis`: Redis client.
- `hashicorp/consul/api`: Runtime configuration source.
- `hashicorp/vault/api`: Secret management source.
- `pressly/goose`: Database migrations.

## Running with Docker Compose

1. Ensure Docker and Docker Compose are installed.
2. Build and run:

```bash
docker-compose up --build
```

3. API is available on `http://localhost:8080`.
4. Stop and remove containers/volumes:

```bash
docker-compose down -v
```

## Configuration Overview

Runtime configuration values are loaded from Consul:

- MySQL connection settings.
- Redis connection settings.
- SMTP settings for verification emails.
- General app settings (domain, http listener port).
- Registration rules (verification expiration, resend limits).
- Registration password policy.

JWT secrets are loaded from Vault:

- `secret/jwt/jwtSecret`
- `secret/jwt/jwtRefreshSecret`

For local development, see `docker-compose.yml` and environment variables used by Vault/Consul clients.

### Logging Configuration

The app supports structured logging with pluggable logger abstraction and Zerolog implementation.

- `LOG_LEVEL`: `debug` | `info` | `warn` | `error` (default: `info`)
- `LOG_FORMAT`: `pretty` | `json`
  - default is `pretty` when `APP_ENV != production`
  - default is `json` when `APP_ENV = production`

Every request gets a correlation ID via `X-Correlation-ID` (auto-generated if missing) and it is included in logs.

### Observability Endpoints

- `GET /metrics` exposes Prometheus metrics.
- `GET /health/live` returns liveness status.
- `GET /health/ready` returns readiness status (MySQL and Redis checks).

The app also includes an OTel-ready middleware skeleton that preserves incoming `traceparent` and injects it into request context/logs for future OpenTelemetry span integration.

## License

This project is licensed under the MIT License. See `LICENSE.txt`.
