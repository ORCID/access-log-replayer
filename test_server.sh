#!/bin/bash
# test_server.sh - Start a simple HTTP server for testing

PORT=8983

# Start a Python HTTP server in the background
python3 -m http.server $PORT &
SERVER_PID=$!
echo "Started test server on port $PORT (PID: $SERVER_PID)"

# Wait a moment for the server to start
sleep 2

echo "You can now run your Go CLI tool to replay requests against http://localhost:$PORT"

echo "Press [Enter] to stop the test server."
read
kill $SERVER_PID
