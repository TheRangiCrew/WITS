#!/bin/bash
set -e

DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER="wits"
DB_NAME="wits"
DB_PASS="WITSTrinity773312*"

echo "Initializing database..."

FILES=("public" "postgis" "awips" "vtec" "mcd")

# Load tables
for sql_file in ${FILES[@]}; do        
    echo "Running $sql_file.sql"
    PGPASSWORD=$DB_PASS psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "../schemas/$sql_file.sql"
done
    
PGPASSWORD=$DB_PASS psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -c "ALTER DATABASE wits SET search_path = public, postgis, awips, vtec, mcd"

# Load data
for sql_file in states offices vtec; do
    echo "Running $sql_file"
    PGPASSWORD=$DB_PASS psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "../data/$sql_file.sql"
done


echo "Done."
