# oc_graylog

---
Tooling to ingest and analyze owncloud.log and audit.log on your local machine using graylog.

## Requirements
- Docker & Compose
- Golang (for building oc_ingest which is faster than the shell-scripts )

## Usage
- Start the graylog stack with `docker-compose up -d`
- Open the graylog web interface at http://127.0.0.1:9000 and login with admin/admin to verify that the stack is running.
- Ingest owncloud.log or audit.log using `ingest_tools/ingenst.sh owncloud.log`
- After ingestion is complete, you can search the logs here http://127.0.0.1:9000/streams/000000000000000000000001/search

### Optional
- Build oc_ingest binary for faster ingestion with  `make oc_ingest` (requires Golang)
- Ingest with `cat owcncloud.log | oc_ingest` 

