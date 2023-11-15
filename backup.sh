# Source environment variables from .env file
if [ -f .env ]; then
    source .env.db
else
    echo "Error: .env.db file not found." >> "$LOG_FILE"
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

# Ensure the backup directory exists
mkdir -p "$BACKUP_DIR"
mkdir -p "$LOG_DIR"

LOG_FILE="$LOG_DIR/backup.log"
touch "$LOG_FILE"

# Generate backup file name with timestamp
BACKUP_FILE="$BACKUP_DIR/backup_$(date +\%Y\%m\%d_\%H\%M\%S).sql"

# Perform backup using pg_dump
PGPASSWORD="$DB_PASSWORD" pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" > "$BACKUP_FILE"

# Check if the backup was successful
if [ $? -eq 0 ]; then
    timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    echo "$timestamp INFO: Backup completed successfully - $BACKUP_FILE" >> "$LOG_FILE"

    # Delete backups older than a week
    find "$BACKUP_DIR" -name 'backup_*.sql' -type f -mtime +7 -delete
else
    timestamp=$(date "+%Y-%m-%d %H:%M:%S")
    echo "$timestamp ERROR: Backup failed. Please check for errors." >> "$LOG_FILE"
    exit 1
fi

# To Run this script (everyday at midnight)
# crontab -e
# 0 0 * * */path/to/backup.sh