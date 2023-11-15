FROM golang:1.20

# Install cron
RUN apt-get update && apt-get -y install cron

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./

# Install app dependencies
RUN go mod download

COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /main

EXPOSE 8000

RUN chmod +x /app/backup.sh

# Apply the cron job
RUN chmod +x /app/cron
RUN crontab /app/cron

# Run
CMD ["/main"]
