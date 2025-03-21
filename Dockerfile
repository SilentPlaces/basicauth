FROM golang:1.23.4 AS builder

# Set the working directory in the container
WORKDIR /app

#installing netcat, i used it in entrypoint.sh (nc -z mysql 3306)
RUN apt-get update && apt-get install -y netcat-openbsd

# Install air and goose
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# set go bin
ENV PATH=$PATH:/go/bin

# download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application's source code
COPY . .

# Copy the entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
# ensure entrypoint.sh is executable
RUN chmod +x /app/entrypoint.sh

# Change working directory to where main.go is placed
WORKDIR /app/cmd/basicauth

# Expose the port your app runs on
EXPOSE 8080

# Use the entrypoint script to run migrations and then start air
CMD ["/bin/sh", "-c", "chmod +x /app/entrypoint.sh && /app/entrypoint.sh"]