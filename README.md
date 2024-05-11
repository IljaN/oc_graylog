# oc_graylog

Tooling to ingest and analyze owncloud.log and audit.log on your local machine using graylog.

## Requirements
- Docker & Compose

## Usage
- Start the graylog stack with `docker-compose up -d`
- Open the graylog web interface at http://127.0.0.1:9000 and login with admin/admin to verify that the stack is running.
- Ingest one or multiple log files with `cat owncloud.log audit.log | docker compose run -T oc_ingest`
- After ingestion is complete, you can search the logs here http://127.0.0.1:9000/streams/000000000000000000000001/search


