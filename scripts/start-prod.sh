#!/bin/bash

COMMAND="./bin/httpserver"
export PORT=443
export REDIRECT_PORT=80

CERT_FILE="/etc/letsencrypt/live/igormichalak.com/fullchain.pem"
KEY_FILE="/etc/letsencrypt/live/igormichalak.com/privkey.pem"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
LOG_FILE="./logs/server_log_$TIMESTAMP.json"
PID_FILE="./logs/server_pid.txt"

mkdir -p ./logs

nohup sudo -E $COMMAND --cert-file=$CERT_FILE --key-file=$KEY_FILE 2>> $LOG_FILE &
PID=$!

echo $PID > $PID_FILE
