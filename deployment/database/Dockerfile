FROM postgis/postgis:17-3.5 AS base

# Install pg_cron
RUN apt-get update && \
    apt-get install -y postgresql-17-cron && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*