FROM postgres:16

WORKDIR /app

# Install cron
RUN apt-get update && apt-get -y install cron

COPY backup.sh /app/backup.sh
COPY restore.sh /app/restore.sh

RUN chmod +x /app/backup.sh
RUN chmod +x /app/restore.sh

# Apply the cron job
RUN chmod +x /app/cron
RUN crontab /app/cron

# Run
CMD ["/main"]
