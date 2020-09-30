#!/bin/sh
# wait-for-postgres.sh

set -e
  
cmd="$1"
  
until PGPASSWORD=$PROVIDER_DB_PASSWORD psql -h "$PROVIDER_DB_HOST" -p "$PROVIDER_DB_PORT" -d "$PROVIDER_DB_NAME" -U "$PROVIDER_DB_USER" -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done
  
>&2 echo "Postgres is up - executing command"
exec $cmd