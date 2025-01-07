#!/bin/bash

PID_FILE="./logs/server_pid.txt"

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo "Stopping server with PID $PID..."
        kill $PID
        if [ $? -eq 0 ]; then
            echo "Server stopped successfully."
            rm -f "$PID_FILE"
        else
            echo "Failed to stop the server."
        fi
    else
        echo "No process found with PID $PID. Cleaning up stale PID file."
        rm -f "$PID_FILE"
    fi
else
    echo "No PID file found."
fi
