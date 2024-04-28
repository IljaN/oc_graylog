#!/usr/bin/bash

# Synchronous owncloud.log ingest
# Usage:
# ingest.sh owncloud.log
while IFS= read -r line; do
    # Parse the field named "time" and covert it to a unix timestamp
    time=$(echo "$line" | jq -r '.time' | xargs -I {} date -d {} +%s)
    # Add a new field named "timestamp" with the unix timestamp to $line without echo
    modified_line=$(echo "$line" | jq --arg time "$time" '. + {version: "1.1", timestamp: $time}')
    # Send to gelf
    echo -n "$modified_line" | nc -w0 -u 127.0.0.1 12201
done < "$1"