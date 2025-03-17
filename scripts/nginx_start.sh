#!/bin/sh

# Handle shutdown gracefully
shutdown() {
  echo "Received shutdown signal, stopping services..."
  nginx -s quit
  kill $(jobs -p) 2>/dev/null
  exit 0
}

# Set up signal handling
trap shutdown SIGTERM SIGINT

# Start Nginx
echo "Starting Nginx..."
nginx

# Start Go application
echo "Starting Go application..."
/usr/local/bin/main &

# Store Go app PID
GO_PID=$!

echo "All services started"

# Wait for the Go application process to exit
wait $GO_PID

# If Go app exits, shut everything down
echo "Go application exited, shutting down container..."
nginx -s quit
exit 0