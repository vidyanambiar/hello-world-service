#!/bin/bash

# Start server
echo "*** Starting API server"
make run & 
export APP_PID=$!
sleep 2

# Run functional tests
make test-functional &
sleep 2

# Shutdown server
kill -s SIGINT -$APP_PID
wait $APP_PID
echo "*** Stopped API server"
