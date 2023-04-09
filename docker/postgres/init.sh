#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER diamondik;
    CREATE DATABASE text_rpg;
    GRANT ALL PRIVILEGES ON DATABASE text_rpg TO diamondik;
EOSQL