#!/usr/bin/bash

# Parallels owncloud.log ingest
# Usage:
#ingest_para.sh owncloud.log
process_line() {
    line=$1

    # Parse the field named "time" and covert it to a unix timestamp
    time=$(echo "$line" | jq -r '.time' | xargs -I {} date -d {} +%s)
    # Add a new field named "timestamp" with the unix timestamp to $line without echo
    modified_line=$(echo "$line" | jq --arg time "$time" '. + {version: "1.1", timestamp: $time}')

    #echo "$modified_line"
    echo -n "$modified_line" | nc -w0 -u 127.0.0.1 12201
}

# Export the function so that it's available to parallel
export -f process_line

# Use GNU Parallel to process the file in parallel with 8 workers
cat "$1" | parallel -j 16 process_line