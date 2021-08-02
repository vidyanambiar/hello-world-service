#!/bin/bash

# Start server
echo "*** Starting API server"
make run & 
export APP_PID=$!
sleep 2

# Run functional tests
make functional-test &
sleep 2

# Shutdown server
kill -s SIGINT -$APP_PID
wait
echo "*** Stopped API server"
