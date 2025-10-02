#!/usr/bin/env bash

# Compile the access-log-replayer program
go build -o access-log-replayer

PORT=8983

# Start a Python HTTP server in the background
python3 custom_http.py &
SERVER_PID=$!

# Wait a moment for the server to start
sleep 2
# Run the access-log-replayer program with test.log
./access-log-replayer -input-file test.log -http_host localhost:8983 > test.debug

kill $SERVER_PID

#git diff --quiet -- server.log
git diff -- server.log

git diff --quiet -- test.debub

# If git diff returns a non-zero exit code, the files have changed
#if [ $? -ne 0 ]; then
#    echo "FATAL: example.properties.j2 or example.properties.yml or example.properties.env.j2 have changed"
#    exit 1
#else
#    echo "No changes in example.properties.j2 or example.properties.yml or example.properties.env.j2"
#fi

