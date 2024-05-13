#!/bin/bash

DB_HOST=postgresdb
DB_DRIVER=postgres
DB_USER=spuser
DB_PASSWORD=SPuser96
DB_NAME=project
DB_PORT=5432


docker run -it -p 1234:5432 \
-e POSTGRES_USER=${DB_USER} \
-e POSTGRES_PASSWORD=${DB_PASSWORD} \
-e POSTGRES_DB=${DB_NAME} \
-e DATABASE_HOST=${DB_HOST} \
--name postgres_a postgres:latest

#psql --host=localhost --port=5432 --username=spuser --dbname=project --password