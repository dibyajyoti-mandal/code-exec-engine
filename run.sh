#!/bin/bash

cleanup() {
    echo "Shutting down all services..."
    kill $GATEWAY_PID 2>/dev/null
    kill $WORKER_PID 2>/dev/null
    kill $FRONTEND_PID 2>/dev/null
    exit
}

trap cleanup SIGINT

echo "Starting Code Execution Engine..."

#Redis
docker start redis
if [ $? -ne 0 ]; then
    echo "Failed to start Redis."
    exit 1
fi

#Gateway Server
cd gateway_server
go run . &
GATEWAY_PID=$!
cd ..

#Start Worker
cd worker
go run . &
WORKER_PID=$!
cd ..

#UI
npm run dev &
FRONTEND_PID=$!

wait