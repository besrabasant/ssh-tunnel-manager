#!/bin/bash

# Ensure the script is run as root
if [ "$(id -u)" -ne 0 ]; then
   echo "This script must be run as root" 1>&2
   exit 1
fi

mv sshtmd /usr/local/bin
mv sshtm /usr/local/bin

# Optional: Set up the application as a system service (for systems using systemd)
cat <<EOF > /etc/systemd/system/sshtmd.service
[Unit]
Description=SSH Tunnnel manager Daemon

[Service]
Type=simple
ExecStart=/usr/local/bin/sshtmd

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload

# Enable and start the service
systemctl enable sshtmd
systemctl start sshtmd