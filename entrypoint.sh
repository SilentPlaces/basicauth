#!/bin/bash
set -e

echo "Waiting for MySQL to be ready..."
until nc -z mysql 3306; do
  sleep 2
done
echo "MySQL is up!"


echo "Running goose migrations..."

goose -dir /app/migrations mysql "user:password@tcp(mysql:3306)/authentication_db" up

echo "Starting air..."
exec air
