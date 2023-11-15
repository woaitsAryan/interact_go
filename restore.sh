#!/bin/bash

# Source database configuration from .env.db file
if [ -f /path/to/.env.db ]; then
    source .env.db
else
    echo "$(date '+%Y-%m-%d %H:%M:%S') ERROR: .env.db file not found." >> "$LOG_FILE"
    exit 1
fi

# Check for required environment variables
required_vars=("DB_HOST" "DB_PORT" "DB_NAME" "DB_USER" "DB_PASSWORD" "BACKUP_DIR" "LOG_DIR")

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        echo "Error: $var is not set. Please check your .env file." >> "$LOG_FILE"
        exit 1
    fi
done

# Find the most recent backup file
MOST_RECENT_BACKUP=$(ls -1t "$BACKUP_DIR" | grep '^backup_' | head -n 1)

if [ -z "$MOST_RECENT_BACKUP" ]; then
    echo "$(date '+%Y-%m-%d %H:%M:%S') ERROR: No backup files found in $BACKUP_DIR." >> "$LOG_FILE"
    exit 1
fi

LOG_FILE="$LOG_DIR/backup.log"
touch "$LOG_FILE"

# Full path to the most recent backup file
BACKUP_FILE="$BACKUP_DIR/$MOST_RECENT_BACKUP"

# Restore database with password
if PGPASSWORD="$DB_PASSWORD" pg_restore -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" < "$BACKUP_FILE"; then
    echo "$(date '+%Y-%m-%d %H:%M:%S') INFO: Database restore completed successfully from $MOST_RECENT_BACKUP." >> "$LOG_FILE"
else
    echo "$(date '+%Y-%m-%d %H:%M:%S') ERROR: Database restore failed. Please check for errors." >> "$LOG_FILE"
    exit 1
fi
