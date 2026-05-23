FROM debian:trixie-slim

RUN apt-get update && apt-get install -y ca-certificates libssl3 && rm -rf /var/lib/apt/lists/*

COPY nerubot /usr/local/bin/
COPY migrations /app/migrations

WORKDIR /app
CMD ["nerubot"]
