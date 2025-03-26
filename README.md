# Basic Authentication in Golang

The project demonstrates basic user authentication using JSON Web Tokens (JWT) for secure, stateless authentication in a Go-based web application.

The application offers the following features:
- User registration and login.
- User E-Mail Verification
- JWT-based authentication to secure endpoints.
- Refresh JWT token
- Simple implementation of the MVC architecture for organizing the codebase.
- A basic setup with Docker and Docker Compose to simplify development and deployment.

## Table of Contents
- [Dependencies](#dependencies)
- [Project Structure](#project-structure)
- [Running the Project with Docker Compose](#running-the-project-with-docker-compose)
- [Main Dependencies](#main-dependencies)
- [Configuration Settings](#configuration-settings)
    - [MySQL Connection Settings](#mysql-connection-settings)
    - [Redis Connection Settings](#redis-connection-settings)
    - [SMTP Server Configuration](#smtp-server-configuration)
    - [General Configuration](#general-configuration)
    - [Registration Password Configuration](#registration-password-configuration)
    - [Vault Configuration](#vault-configuration)
- [License](#license)

## Dependencies

This project has the following dependencies:

- **Golang 1.x+**: Required for building and running the Go application.
- **Docker**: Used for containerizing the application and its dependencies.
- **Docker Compose**: Manages the multi-container setup for the development environment (e.g., application container, database, etc.).
- **JWT (JSON Web Tokens)**: Used to implement secure, stateless authentication for API endpoints.

### Project Structure:
- **cmd**: Contains the main application entry point.
- **migrations**: Database migration files (if applicable).
- **internal**: Core application logic (e.g., authentication, user management).
- **pkg**: Utility packages shared across the application.
- **Dockerfile**: Defines the application container.
- **docker-compose.yml**: Configuration for setting up the development environment.
- **entrypoint.sh**: Script to initialize the app inside the container.
- **.env**: Environment variable configuration for local development.
- **basic-auth.postman_collection.json**: Postman collection for this project, which provides ready-to-use API tests and requests.

## Running the Project with Docker Compose

This project uses Docker Compose to set up the development environment. Follow these steps to get it running:

1. Ensure you have Docker and Docker Compose installed on your machine.
2. Clone the repository and navigate to the project folder:

```bash
git clone <repository_url>
cd basicauth
```

3. Build and run the development environment using Docker Compose.

```bash
docker-compose up --build
```

4. Once the container is up, the application will be available at `http://localhost:8080` (or another port defined in the Docker Compose file).

5. If you need to stop the application, simply use:

```bash
docker-compose down -v
```

**Important:**
- The Vault environment variables in the Docker Compose file are for development purposes only. Do **not** use this approach in a production environment. Consider using a secure secret management solution like HashiCorp Vault or AWS Secrets Manager for production deployments.

## Main Dependencies

### 1. **httprouter**
Fast HTTP router for defining routes and handling requests.

### 2. **go-redis**
Go client for interacting with Redis, used for caching and fast data retrieval.

### 3. **go-sql-driver/mysql**
MySQL driver for Go, enabling interaction with MySQL databases.

### 4. **consul/api**
Go client for HashiCorp Consul, used for configuration management and service discovery.

### 5. **vault/api**
Go client for HashiCorp Vault, used for securely managing secrets and sensitive data.

### 6. **wire**
Dependency injection framework for managing dependencies in Go applications.

### 7. **golang-jwt/jwt**
Library for creating and verifying JWTs (JSON Web Tokens), used for authentication.

### 8. **goose**
Database migration tool for Go, used for managing schema changes.

## Configuration Settings

This section explains the configuration options stored in Consul for various services used in this project.

### MySQL Connection Settings
These settings are used for connecting to a MySQL database:
- **host**: The hostname of the MySQL server. Default is `mysql`.
- **user**: The MySQL user for authentication. Default is `user`.
- **port**: The port number for the MySQL server. Default is `3306`.
- **password**: The password for the MySQL user. Default is `password`.
- **db**: The name of the database used for authentication. Default is `authentication_db`.
- **maxLifeTime**: The maximum lifetime of a connection in seconds. Default is `180`.
- **idleConnections**: The number of idle connections to maintain. Default is `10`.
- **maxOpenConnections**: The maximum number of open connections. Default is `10`.

### Redis Connection Settings
These settings are used for connecting to a Redis server:
- **host**: The hostname of the Redis server. Default is `redis`.
- **port**: The port number for the Redis server. Default is `6379`.
- **password**: The password for the Redis server. Default is an empty string `''`.

### SMTP Server Configuration
These settings are used to configure the SMTP server for email sending:
- **host**: The hostname of the SMTP server. Default is `mailhog`.
- **port**: The port number for the SMTP server. Default is `1025`.
- **username**: The SMTP username for authentication. Default is an empty string `''`.
- **password**: The SMTP password for authentication. Default is an empty string `''`.

### General Configuration
This section covers general configuration options:
- **domain**: The domain of the application, default is `localhost`.
- **httpListenerPort**: The port where the application listens for HTTP requests. Default is `8080`.
- **register/maxVerificationMailInCountInDay**: The maximum number of verification emails sent per day for registration. Default is `5`.
- **register/mailVerificationTimeInSeconds**: The expiration time for the email verification link in seconds. Default is `6000`.
- **register/hostVerificationMailAddress**: The email address used as the sender for verification emails. Default is `armin@testlocalhost.com`.
- **register/verificationMailText**: The content of the verification email sent to users. This is an HTML email template with a link to verify the email address.

### Registration Password Configuration
This section defines password requirements for user registration:
- **minLength**: The minimum length of the password. Default is `8`.
- **requireUpper**: Whether the password must contain an uppercase letter. Default is `false`.
- **requireLower**: Whether the password must contain a lowercase letter. Default is `false`.
- **requireNumber**: Whether the password must contain a number. Default is `false`.
- **requireSpecial**: Whether the password must contain a special character. Default is `false`.

### Vault Configuration
This section covers the JWT secrets stored in Vault for authentication purposes:

- **secret/jwt/jwtSecret**: The secret key used to sign JWT tokens. This key is crucial for the generation of valid tokens and should be kept secret to ensure the integrity of authentication in the application. Default value is `jwt_secret_value`.

- **secret/jwt/jwtRefreshSecret**: The secret key used to sign JWT refresh tokens. This key is used for generating refresh tokens that allow users to obtain a new JWT when the current one expires. Default value is `jwt_refresh_token_value`.

These JWT secrets are stored securely in Vault and used for generating tokens during authentication and session management.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.txt) file for details.
