#!/bin/bash
set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER="wits"
DB_NAME="wits"
DB_PASS="WITSTrinity773312*"

echo "Initializing database..."

FILES=("cwa.sql" "counties.sql" "zones.sql" "marinezones.sql" "firezones.sql")

for sql_file in ${FILES[@]}; do        
    echo "  Running $sql_file.sql"
    PGPASSWORD=$DB_PASS psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$sql_file"
done

echo "Done."
