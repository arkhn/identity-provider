#!/bin/sh

TMP_FILE=/tmp/usersdump
touch $TMP_FILE
trap "rm $TMP_FILE" EXIT

echo "Copy users from source..."
PGPASSWORD=${SOURCE_DB_PASSWORD} psql -h ${SOURCE_DB_HOST} -p ${SOURCE_DB_PORT} -U ${SOURCE_DB_USER} -d ${SOURCE_DB} -c "COPY \"${SOURCE_SCHEMA}\".\"${SOURCE_TABLE}\" (email, name, password) TO '$TMP_FILE'"

echo "Write users to target..."
PGPASSWORD=${PROVIDER_DB_PASSWORD} psql -h ${PROVIDER_DB_HOST} -p ${PROVIDER_DB_PORT} -U ${PROVIDER_DB_USER} -d ${PROVIDER_DB_NAME} -c "COPY users (email, name, password) FROM '$TMP_FILE'"
