#!/bin/bash

DB_HOST=localhost
DB_DRIVER=postgres
DB_USER=spuser
DB_PASSWORD=SPuser96
DB_NAME=project
DB_PORT=1234

POSTGRES_USER=spuser \
POSTGRES_PASSWORD=SPuser96 \
POSTGRES_DB=project \
DATABASE_HOST=localhost \
DATABASE_PORT=1234 \
CONFIG_PATH=./config/prod.yaml go run cmd/url-shortener/main.go