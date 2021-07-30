#!/bin/bash

echo "*** Starting API server"
make run &
export APP_PID=$!
make test-functional
echo "*** Stopping API server"
kill $APP_PID