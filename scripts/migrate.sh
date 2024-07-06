#!/bin/sh
echo "Starting migrations"
export PGPASSWORD=$DB_PASSWORD


for filename in ./scripts/db/migrations/*; do
    if [ -f "$filename" ]; then
        echo "Running migration $filename"
        psql -h $DB_HOST -d $DB_NAME -U $DB_USERNAME -f $filename
        if [ $? -eq 0 ]; then
            echo "Migration successful $filename"
        else
            echo "Error running migration."
            exit 1
        fi
    fi
done
echo "Migration completed."