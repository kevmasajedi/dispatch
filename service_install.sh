#!/bin/bash

SERVICE_NAME="dispatch" 
BINARY_NAME="dispatch"     
WORKING_DIR=$(pwd)       
USER=$(whoami)       

# Compile the program
go build -o $BINARY_NAME main.go

# Generate the systemd service file
SERVICE_FILE_CONTENT="[Unit]
Description=$SERVICE_NAME Service
After=network.target

[Service]
ExecStart=$WORKING_DIR/$BINARY_NAME
WorkingDirectory=$WORKING_DIR
Restart=always
RestartSec=3
User=$USER
Group=$USER

[Install]
WantedBy=multi-user.target
"

echo "Creating service file for $SERVICE_NAME at /etc/systemd/system/$SERVICE_NAME.service"
echo "$SERVICE_FILE_CONTENT" |  tee /etc/systemd/system/$SERVICE_NAME.service > /dev/null

systemctl daemon-reload

systemctl enable $SERVICE_NAME.service

systemctl start $SERVICE_NAME.service

echo "Service status:"
systemctl status $SERVICE_NAME.service

# Output message
echo "Service $SERVICE_NAME created, started, and enabled at boot successfully."
