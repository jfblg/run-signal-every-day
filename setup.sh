#!/bin/bash

go build signalkeepalive.go
echo "Signal Keepalive built successfully."

mkdir -p "$HOME/signal-keepalive"
echo "Directory $HOME/signal-keepalive created successfully."

cp signalkeepalive "$HOME/signal-keepalive/"
echo "File keepalive-signal copied to $HOME/signal-keepalive successfully."

chmod +x "$HOME/signal-keepalive/signalkeepalive"
echo "File signalkeepalive made executable successfully."

cp net.jfblg.signalkeepalive.plist "$HOME/Library/LaunchAgents/"
echo "File net.jfblg.signalkeepalive.plist copied to $HOME/Library/LaunchAgents successfully."

launchctl load "$HOME/Library/LaunchAgents/net.jfblg.signalkeepalive.plist"
echo "Signal Keepalive service loaded successfully."

launchctl start net.jfblg.signalkeepalive
echo "Signal Keepalive service started successfully."

echo "Setup completed successfully."
