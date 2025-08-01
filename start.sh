#!/bin/sh

set -e 

echo "starting database migration" 
/app/migrate -path /app/migration -database "$DB_SOURCE" -verbose up

echo "starting app:"
exec "$@"