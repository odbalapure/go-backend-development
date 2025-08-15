#! /bin/sh

# Exit script after getting a non-zero exit code
set -e

echo "run db migration"

source /app/app.env

# /app/migrate is the binary
# -path is the path to the migration files
# -database is the database connection string
# -verbose is to print the migration progress
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "start the app"

# Takes all the arguments passed to the script and runs them
exec "$@"
