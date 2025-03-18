#!/bin/sh
set -e

# Wait for the database to be ready
echo "Waiting for PostgreSQL to be ready..."
RETRIES=10
until pg_isready -h db -U app -d appdb || [ $RETRIES -eq 0 ]; do
  echo "Waiting for PostgreSQL server, $((RETRIES--)) remaining attempts..."
  sleep 2
done

# Start the application
echo "Starting application..."
exec "$@"