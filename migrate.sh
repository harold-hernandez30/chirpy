#! /usr/bin/bash -e

MIGRATIONS_DIR=sql/schema
DATABASE_URL="postgres://postgres:postgres@localhost:5432/chirpy"

case $1 in
    up)
        goose -dir $MIGRATIONS_DIR postgres $DATABASE_URL up
    ;;
    down)
        goose -dir $MIGRATIONS_DIR postgres $DATABASE_URL down
    ;;
    down-to-zero)
        goose -dir $MIGRATIONS_DIR postgres $DATABASE_URL down-to 0
    ;;
    *)
        echo "argument $1 - is not supported"
    ;;
esac