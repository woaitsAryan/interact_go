FROM postgres:16

WORKDIR /app

# Install cron
RUN apt-get update && apt-get -y install cron

COPY backup.sh /app/backup.sh
COPY remote_backup.sh /app/remote_backup.sh
COPY restore.sh /app/restore.sh

COPY db.cron /app/db.cron

COPY .env.db /app/.env.db
COPY .env.db.remote /app/.env.db.remote

RUN chmod +x /app/backup.sh
RUN chmod +x /app/remote_backup.sh
RUN chmod +x /app/restore.sh

# Apply the cron job
RUN chmod +x /app/db.cron
RUN crontab /app/db.cron
