#! /usr/bin/bash -e

MIGRATIONS_DIR=sql/schema
DATABASE_URL="postgres://postgres:postgres@localhost:5432/chirpy"


goose -dir $MIGRATIONS_DIR postgres $DATABASE_URL $1